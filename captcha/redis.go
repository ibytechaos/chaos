/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2022 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: redis.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2022-11-20
 * | Description: redis.go
 * |-----------------------------------------------------------
 */

package captcha

import (
	"github.com/ibytechaos/chaos/client/redis"
	"time"
)

func NewRedisStore() *RedisStore {
	return &RedisStore{
		Expiration: time.Second * 180,
		PreKey:     "CAPTCHA_",
	}
}

// RedisStore redis存储
type RedisStore struct {
	// Expiration 验证码过期时间
	Expiration time.Duration
	// PreKey 前缀
	PreKey string
}

// Set 设置
func (o *RedisStore) Set(id string, value string) error {
	return redis.GetInstance().AsyncSet(&redis.Data{
		Key:        o.Key(id),
		Value:      value,
		Expiration: o.Expiration,
	})
}

// Get 获取
func (o *RedisStore) Get(key string, clear bool) string {
	val, err := redis.GetInstance().Get(o.Key(key))
	if err != nil {
		return ""
	}
	if clear {
		err = redis.GetInstance().Del(o.Key(key))
	}
	return val
}

// Verify 验证
func (o *RedisStore) Verify(id, answer string, clear bool) bool {
	key := o.PreKey + id
	v := o.Get(key, clear)
	return v == answer
}

// Key 根据id生成key
func (o *RedisStore) Key(key string) string {
	return o.PreKey + key
}
