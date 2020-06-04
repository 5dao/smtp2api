//Package server svr
package server

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/5dao/golibs/log"
	"github.com/gin-gonic/gin"
)

const tokenSolt = "shishuilingqingshanyebian:"

// Server send mail
type Server struct {
	Cfg    *Config
	GinEg  *gin.Engine
	Locker *sync.Mutex

	today          time.Time
	nextIndex      int
	todaySendCount int
}

// NewServer make new
func NewServer(cfg *Config) (svr *Server, err error) {
	if len(cfg.Accounts) < 1 {
		return nil, fmt.Errorf("atleast one account: %v", cfg.Accounts)
	}

	for _, account := range cfg.Accounts {
		hostPort := strings.Split(account.SMTP, ":")
		if len(hostPort) != 2 {
			return nil, fmt.Errorf("smtp addr err: %s", account.SMTP)
		}
		port, err := strconv.ParseInt(hostPort[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("smtp addr err: %v,%s", err, account.SMTP)
		}
		account.smtpHost = hostPort[0]
		account.smtpPort = int(port)
	}

	svr = &Server{
		Cfg:    cfg,
		GinEg:  gin.Default(),
		Locker: new(sync.Mutex),

		today: time.Now(),
	}

	svr.InitRouter()

	return
}

//Start  server
func (svr *Server) Start() {
	// http server
	go svr.GinEg.Run(svr.Cfg.Listen)
	log.Println("svr: start ok!")
}

//InitRouter InitRouter
func (svr *Server) InitRouter() {
	svr.GinEg.Use(Cors())
	svr.GinEg.POST(svr.Cfg.BasePath+"/mailto", svr.HandleMailTo)
	svr.GinEg.GET(svr.Cfg.BasePath+"/test", func(c *gin.Context) {
		c.JSON(200, &JSONResult{
			Code: 0,
			Msg:  "ok,RequestURI:" + c.Request.RequestURI,
		})

		return
	})
}

// Cors Cors
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

//JSONResult JSONResult
type JSONResult struct {
	Code int    `json:"code"` //0 true,
	Msg  string `json:"msg"`
}

//HandleMailTo MailTo
func (svr *Server) HandleMailTo(c *gin.Context) {
	var to, cc []string
	var subject, body string

	to = svr.Cfg.Addrs

	rs := &JSONResult{
		Code: -1,
		Msg:  "",
	}

	//token check
	now := time.Now()

	var err error

	var formToken, token, formTimestamp string
	var formTo []string
	var keyExist bool
	var timestrap, timestrapLB int64

	if formToken, keyExist = c.GetPostForm("token"); !keyExist {
		rs.Code = 100
		rs.Msg = "token need"
		goto RS
	}
	token, err = AesDecryptToken(formToken, svr.Cfg.TokenKey)
	if err != nil {
		rs.Code = 101
		rs.Msg = "token err"
		goto RS
	}

	if formTimestamp, keyExist = c.GetPostForm("t"); !keyExist {
		rs.Code = 110
		rs.Msg = "timestrap need"
		goto RS
	}
	timestrap, err = strconv.ParseInt(formTimestamp, 10, 64)
	if err != nil {
		rs.Code = 111
		rs.Msg = "timestrap err"
		goto RS
	}
	timestrapLB = timestrap - now.Unix()
	if timestrapLB > 30 || timestrapLB < -30 {
		rs.Code = 112
		rs.Msg = "timestrap over err"
		goto RS
	}

	if token != (tokenSolt + formTimestamp) {
		rs.Code = 300
		rs.Msg = "token err"
		goto RS
	}

	if formTo, keyExist = c.GetPostFormArray("to"); keyExist {
		to = append(to, formTo...)
	}

	cc, _ = c.GetPostFormArray("cc")
	subject, _ = c.GetPostForm("subject")
	body, _ = c.GetPostForm("body")

	if err := svr.SendTo(to, cc, subject, body); err != nil {
		rs.Code = 101
		rs.Msg = fmt.Sprintf("send server error: %v", err)
		goto RS
	}

	rs.Code = 0
	rs.Msg = "send ok"

RS:
	c.JSON(200, rs)
}

func (svr *Server) getBestAccount() *MailAccount {
	svr.Locker.Lock()
	defer svr.Locker.Unlock()

	now := time.Now()

	m1Time := now.Add(-1 * time.Minute)

	if svr.today.Day() != now.Day() {
		//reset new day
		for _, account := range svr.Cfg.Accounts {
			account.TodayCount = 0
			account.LastSendTime = now.Add(-10 * time.Minute)

			svr.todaySendCount = 0
		}
		svr.today = now
	}

	if svr.nextIndex >= len(svr.Cfg.Accounts) {
		svr.nextIndex = 0
	}

	cAccount := svr.Cfg.Accounts[svr.nextIndex]
	log.Printf("SendTo: from: %d,%s", svr.nextIndex, cAccount.User)
	if cAccount.LastSendTime.Before(m1Time) && (cAccount.TodayCount < cAccount.Max) {
		svr.nextIndex = svr.nextIndex + 1
		cAccount.LastSendTime = now
		return cAccount
	}

	var loopI = svr.nextIndex + 1
	for loopI != svr.nextIndex {
		cAccount = svr.Cfg.Accounts[loopI]
		if cAccount.LastSendTime.Before(m1Time) && (cAccount.TodayCount < cAccount.Max) {
			svr.nextIndex = loopI + 1
			cAccount.LastSendTime = now
			return cAccount
		}

		loopI = loopI + 1
		if loopI >= len(svr.Cfg.Accounts) {
			loopI = 0
		}
	}
	return nil
}

// SendTo mail
func (svr *Server) SendTo(to []string, cc []string, subject string, body string) (err error) {
	// for i := 0; i < 3; i++ {
	log.Printf("SendTo: index: %d, to: %v, cc: %v, t: %s, msg: %s", svr.todaySendCount, to, cc, subject, body)

	//get one account per send,so 30 second ,not 1 minute
	account := svr.getBestAccount()

	if account == nil {
		err = fmt.Errorf("no account to use: %s", "")
		log.Printf("SendTo: index: %d, err: %v", svr.todaySendCount, err)
		return
	}

	if err = SendMailByAccount(account, to, cc, svr.Cfg.SubjectPrefx+": "+subject, body); err != nil {
		err = fmt.Errorf("do err: %v", err)
		log.Printf("SendTo: index: %d, err: %v", svr.todaySendCount, err)
		// time.Sleep(time.Second * 60)
		// continue
	}

	log.Printf("SendTo: index: %d,ok", svr.todaySendCount)

	svr.todaySendCount = svr.todaySendCount + 1

	// saveF := fmt.Sprintf("logs/%s-%d.md", svr.today.Format("01-02-1504"), svr.todaySendCount)
	// txt = subject + `
	// ` + body
	// err = ioutil.WriteFile(saveF, []byte(txt), os.ModePerm)
	return
	// }
	// return
}
