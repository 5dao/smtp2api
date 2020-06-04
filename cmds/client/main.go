//Package main
package main

import (
	"flag"
	"log"
	"time"

	"github.com/5dao/smtp2api/postmail"
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
	rs, err := postmail.PostMail(time.Now(), "shishuilingqingshanyebian:", apiURL, subject, body, key)
	log.Println(rs, err)
}