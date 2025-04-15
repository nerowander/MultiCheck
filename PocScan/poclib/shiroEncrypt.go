package poclib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	uuid "github.com/satori/go.uuid"
	"io"
)

var (
	// 这个checkcontent是用来检测shiro key是否正确的，用的是SimplePrincipalCollection类
	CheckContent = "rO0ABXNyADJvcmcuYXBhY2hlLnNoaXJvLnN1YmplY3QuU2ltcGxlUHJpbmNpcGFsQ29sbGVjdGlvbqh/WCXGowhKAwABTAAPcmVhbG1QcmluY2lwYWxzdAAPTGphdmEvdXRpbC9NYXA7eHBwdwEAeA=="
	Content, _   = base64.StdEncoding.DecodeString(CheckContent)
)

func Padding(plainText []byte, blockSize int) []byte {
	n := (blockSize - len(plainText)%blockSize)
	temp := bytes.Repeat([]byte{byte(n)}, n)
	plainText = append(plainText, temp...)
	return plainText
}

func GetShiroCookie(key, mode string) string {
	if mode == "gcm" {
		return AES_GCM_Encrypt(key)
	} else {
		//cbc
		return AES_CBC_Encrypt(key)
	}
}

//AES CBC加密
func AES_CBC_Encrypt(shirokey string) string {
	key, err := base64.StdEncoding.DecodeString(shirokey)
	if err != nil {
		return ""
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	Content = Padding(Content, block.BlockSize())
	iv := uuid.NewV4().Bytes()
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cipherText := make([]byte, len(Content))
	blockMode.CryptBlocks(cipherText, Content)
	return base64.StdEncoding.EncodeToString(append(iv[:], cipherText[:]...))
}

//AES GCM 加密
func AES_GCM_Encrypt(shirokey string) string {
	key, err := base64.StdEncoding.DecodeString(shirokey)
	if err != nil {
		return ""
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}
	nonce := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return ""
	}
	aesgcm, _ := cipher.NewGCMWithNonceSize(block, 16)
	ciphertext := aesgcm.Seal(nil, nonce, Content, nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...))
}
