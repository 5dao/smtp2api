//Package server aes
package server

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
)

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

// AesDecryptToken AesDecryptToken
func AesDecryptToken(token, passphrase string) (string, error) {
	ct, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	key := passphrase
	data := []byte(passphrase)
	iv := md5.Sum(data)
	iva := iv[:]
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, iva)
	dst := make([]byte, len(ct))

	cfb.XORKeyStream(dst, ct)
	return string(dst), nil
}
