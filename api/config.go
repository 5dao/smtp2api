//Package api config
package api

import "time"

//Config api config
type Config struct {
	Listen   string // :27899
	BasePath string // api = http://{{listen}}/{{BasePath}}/mailto

	Accounts     []*MailAccount
	SubjectPrefx string

	Addrs []string

	TokenKey string
}

// MailAccount smtp account
type MailAccount struct {
	SMTP     string //smtp.mail.com:25
	User     string
	Password string
	Max      int //max send per day

	//
	smtpHost string
	smtpPort int

	TodayCount   int //
	LastSendTime time.Time
}
