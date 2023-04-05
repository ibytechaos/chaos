/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2022 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: security.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2022-11-20
 * | Description: security.go
 * |-----------------------------------------------------------
 */

package utils

import (
	"github.com/ibytechaos/chaos/g"
	"github.com/wumansgy/goEncrypt/aes"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	token = getSecret()
}

var token = ""

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptVerify 对比明文密码和数据库的哈希值
func BcryptVerify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Token 令牌
func Token(t string) string {
	t = strings.ToLower(t)
	t = getTokenFromFile(t)
	if t == "" {
		return t
	}
	return Decrypt(getSecret(), t)
}

// SaveToken 保存令牌
func SaveToken(t string, token string) {
	token = Encrypt(getSecret(), token)
	saveTokenToFile(t, token)
}

func getSecret() string {
	if token == "" {
		token = getTokenFromFile("secret")
		if token == "" {
			token = Random(32)
			saveTokenToFile("secret", token)
		}
	}
	return token
}

func saveTokenToFile(t string, token string) {
	if _, err := os.Stat(g.SecurityDir); os.IsNotExist(err) {
		_ = os.Mkdir(g.SecurityDir, 0755)
	}
	if _, err := os.Stat(g.SecurityDir + "/" + t); os.IsNotExist(err) {
		_ = os.Mkdir(g.SecurityDir+"/"+t, 0755)
	}
	files, err := os.ReadDir(g.SecurityDir + "/" + t)
	if len(files) == 0 || err != nil {
		_ = os.WriteFile(g.SecurityDir+"/"+t+"/"+Random(32), []byte(token), 0600)
		return
	}
	if len(files) > 0 {
		_ = os.WriteFile(g.SecurityDir+"/"+t+"/"+files[0].Name(), []byte(token), 0600)
	}
}

func getTokenFromFile(t string) string {
	files, err := os.ReadDir(g.SecurityDir + "/" + t)
	if len(files) == 0 || err != nil {
		return ""
	}
	for _, f := range files {
		data, err := os.ReadFile(g.SecurityDir + "/" + t + "/" + f.Name())
		if err != nil {
			return ""
		}
		return string(data)
	}
	return ""
}

// Encrypt 加密
func Encrypt(key string, token string) string {
	hex, _ := aes.AesCbcEncryptHex([]byte(token), []byte(key), nil)
	return hex
}

// Decrypt 解密
func Decrypt(key string, token string) string {
	dec, _ := aes.AesCbcDecryptByHex(token, []byte(key), nil)
	return string(dec)
}

// Random 随机字符串
func Random(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
