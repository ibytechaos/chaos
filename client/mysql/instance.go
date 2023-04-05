/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2022 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: instance.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2022-11-27
 * | Description: instance.go
 * |-----------------------------------------------------------
 */

package mysql

import (
	"fmt"
	"github.com/ibytechaos/chaos/config"
	"github.com/ibytechaos/chaos/utils"
	"golang.org/x/exp/slog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

// MySQL 单例
type MySQL struct {
	DB *gorm.DB
}

var (
	instance *MySQL
	once     sync.Once
)

// GetInstance 获取配置实例
func GetInstance() *MySQL {
	once.Do(func() {
		instance = &MySQL{}
		instance.init()
	})
	return instance
}

func (o *MySQL) init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.GetInstance().Client.MySQL.Username,
		utils.Token(config.GetInstance().Client.MySQL.Name),
		config.GetInstance().Client.MySQL.Host,
		config.GetInstance().Client.MySQL.Port,
		config.GetInstance().Client.MySQL.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("mysql init error", err)
	}
	o.DB = db
}
