package api

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// ma la ge bi aliyun.com ,so bi

// SendMailByAccount by account
func SendMailByAccount(account *MailAccount, to []string, cc []string, subject string, body string) error {
	// fmt.Println(account.smtpHost, account.smtpPort, account.User, account.Password, to, cc, subject, body)
	// return nil

	return SendMail(account.smtpHost, account.smtpPort, account.User, account.Password, to, cc, subject, body)
}

// SendMail to []addr
func SendMail(host string, port int, account string, pwd string, to []string, cc []string, subject string, body string) (err error) {
	defer func() {
		if rev := recover(); rev != nil {
			err = fmt.Errorf("SendMail: recover: %v,%s,%s,%s", rev, account, subject, body)
		}
	}()

	msg := gomail.NewMessage()
	msg.SetHeader("From", account)
	msg.SetHeader("To", to...) //  see all to addrs in receive addr bar
	for _, ccAddr := range cc {
		msg.SetAddressHeader("Cc", ccAddr, ccAddr) //copy send, see all cc addrs in copy send addr bar
	}
	msg.SetHeader("Subject", subject) //subject
	msg.SetBody("text/html", body)    //body
	//msg.Attach("/home/Alex/lolcat.jpg")//attach files

	cnn := gomail.NewDialer(host, 25, account, pwd)
	//cnn.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return cnn.DialAndSend(msg)
}
