package postmail

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

//JSONResult JSONResult
type JSONResult struct {
	Code int    `json:"code"` //0 true,
	Msg  string `json:"msg"`
}

//PostMail PostMail
//tokenSolt and key
func PostMail(t time.Time, tokenSolt, apiURL, subject, body, key string, tos []string) (rs *JSONResult, err error) {
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

	for _, to := range tos {
		values.Add("to", to)
	}

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
