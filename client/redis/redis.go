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

package redis

import (
	"context"
	"github.com/ibytechaos/chaos/config"
	"github.com/ibytechaos/chaos/g"
	"github.com/ibytechaos/chaos/utils"
	"github.com/panjf2000/ants/v2"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
	"time"
)

// Redis 单例
type Redis struct {
	pool  *ants.PoolWithFunc
	read  redis.UniversalClient
	write redis.UniversalClient
}

var (
	instance *Redis
	once     sync.Once
)

// GetInstance 获取配置实例
func GetInstance() *Redis {
	once.Do(func() {
		instance = &Redis{}
		instance.init()
	})
	return instance
}

func (o *Redis) init() {
	if config.GetInstance().Client.Redis == nil {
		return
	}
	db := 0
	if config.GetInstance().Client.Redis.Database != "" {
		if v, err := strconv.ParseInt(config.GetInstance().Client.Redis.Database, g.DefaultBase, g.DefaultBitSize); err == nil {
			db = int(v)
		}
	}
	rOption := &redis.UniversalOptions{
		Addrs:       config.GetInstance().Client.Redis.Hosts,
		Password:    utils.Token(config.GetInstance().Client.Redis.Token),
		DB:          db,
		ReadTimeout: time.Duration(config.GetInstance().Client.Redis.Timeout) * time.Millisecond,
		MasterName:  config.GetInstance().Client.Redis.MasterName,
	}
	wOption := &redis.UniversalOptions{
		Addrs:       config.GetInstance().Client.Redis.Hosts,
		Password:    utils.Token(config.GetInstance().Client.Redis.Token),
		DB:          db,
		ReadTimeout: time.Duration(config.GetInstance().Client.Redis.Timeout+1000) * time.Millisecond,
		MasterName:  config.GetInstance().Client.Redis.MasterName,
	}
	o.read = redis.NewUniversalClient(rOption)
	o.write = redis.NewUniversalClient(wOption)
	o.pool = o.NewPool()
}

// Data 数据
type Data struct {
	// Key 键
	Key string
	// Value 值
	Value interface{}
	// Expiration 过期时间
	Expiration time.Duration
	// Compress 是否压缩
	Compress bool
}

// NewPool 工作池
func (o *Redis) NewPool() *ants.PoolWithFunc {
	pool, _ := ants.NewPoolWithFunc(g.CpuNumber, func(i interface{}) {
		if data, ok := i.(*Data); ok {
			if data.Compress {
				data.Value = utils.Compress(data.Value.([]byte))
			}
			o.write.Set(context.Background(), data.Key, data.Value, data.Expiration)
		}
	}, ants.WithNonblocking(true))
	return pool
}
