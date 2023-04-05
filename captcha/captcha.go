/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2022 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: captcha.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2022-11-20
 * | Description: captcha.go
 * |-----------------------------------------------------------
 */

package captcha

import (
	"github.com/ibytechaos/chaos/config"
	"github.com/mojocn/base64Captcha"
	"sync"
)

// CR 验证码
type CR struct {
	// id 验证码ID
	Id string `json:"id"`
	// B64 b64s 验证码图片
	B64 string `json:"b64"`
	// Value 验证码值
	Value string `json:"value"`
	// Clear 是否清除
	Clear bool `json:"clear"`
	// Length 验证码长度
	Length int `json:"length"`
}

// Captcha 单例
type Captcha struct {
	captcha *base64Captcha.Captcha
	store   base64Captcha.Store
}

var (
	instance *Captcha
	once     sync.Once
)

// GetInstance 获取配置实例
func GetInstance() *Captcha {
	once.Do(func() {
		instance = &Captcha{}
		instance.init()
	})
	return instance
}

func (o *Captcha) init() {
	if config.GetInstance().Client.Redis == nil {
		o.store = base64Captcha.DefaultMemStore
	} else {
		o.store = NewRedisStore()
	}
	// 生成默认数字的driver
	o.captcha = base64Captcha.NewCaptcha(
		base64Captcha.NewDriverDigit(
			config.GetInstance().App.Captcha.Height,
			config.GetInstance().App.Captcha.Width,
			config.GetInstance().App.Captcha.Length,
			0.7,
			80),
		o.store)
}

// Generate 生成验证码
func (o *Captcha) Generate() (*CR, error) {
	id, b64, err := o.captcha.Generate()
	return &CR{Id: id, B64: b64}, err
}

// Verify 验证验证码
func (o *Captcha) Verify(r *CR) bool {
	return o.store.Verify(r.Id, r.Value, r.Clear)
}
