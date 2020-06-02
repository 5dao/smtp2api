package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var apiURL, subject, body, key string

func init() {
	flag.StringVar(&apiURL, "u", "http://127.0.0.1/api/mailto", "-u ")
	flag.StringVar(&key, "k", "", "-k token key")
	flag.StringVar(&subject, "s", "email subject", "-s subject")
	flag.StringVar(&body, "m", "email body", "-m body")
	flag.Parse()
}

func main() {
	rs, err := PostMail(time.Now(), "shishuilingqingshanyebian:", apiURL, subject, body, key)
	log.Println(rs, err)
}

//
//
//

//JSONResult JSONResult
type JSONResult struct {
	Code int    `json:"code"` //0 true,
	Msg  string `json:"msg"`
}

//PostMail PostMail
//tokenSolt and key
func PostMail(t time.Time, tokenSolt, apiURL, subject, body, key string) (rs *JSONResult, err error) {
	// defer func() {
	// 	if rev := recover(); rev != nil {
	// 		err = fmt.Errorf("PostMail rev: %v", rev)
	// 	}
	// }()

	values := make(url.Values)

	timestampStr := strconv.FormatInt(t.Unix(), 10)
	tokenStr := tokenSolt + timestampStr

	var token string
	token, err = AesEncryptToken(tokenStr, key)
	if err != nil {
		return
	}

	values.Add("token", token)
	values.Add("t", timestampStr)
	values.Add("subject", subject)
	values.Add("body", body)

	var resp *http.Response
	resp, err = http.PostForm(apiURL, values)
	if err != nil {
		return
	}
	var resBody []byte
	resBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	rs = &JSONResult{}
	err = json.Unmarshal(resBody, rs)
	if err != nil {
		return nil, fmt.Errorf("response not json,%s,%v", string(resBody), err)
	}
	return rs, nil
}

//AesEncryptToken AesEncryptToken
func AesEncryptToken(plainText, passphrase string) (string, error) {
	key := passphrase
	data := []byte(passphrase)
	iv := md5.Sum(data)
	iva := iv[:]
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	//pad := __PKCS7Padding([]byte(text), block.BlockSize())
	byteText := []byte(plainText)
	cfb := cipher.NewCFBEncrypter(block, iva)
	encrypted := make([]byte, len(byteText))
	cfb.XORKeyStream(encrypted, byteText)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}
