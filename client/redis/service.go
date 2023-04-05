/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2022 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: service.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2022-11-20
 * | Description: service.go
 * |-----------------------------------------------------------
 */

package redis

import (
	"context"
	"time"
)

// Get 获取
func (o *Redis) Get(key string) (string, error) {
	return o.read.Get(context.Background(), key).Result()
}

// Set 设置
func (o *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return o.write.Set(context.Background(), key, value, expiration).Err()
}

// Del 删除
func (o *Redis) Del(key string) error {
	return o.write.Del(context.Background(), key).Err()
}

// AsyncSet 异步设置
func (o *Redis) AsyncSet(data *Data) error {
	return o.pool.Invoke(data)
}
