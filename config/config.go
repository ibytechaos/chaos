/*
 * |-----------------------------------------------------------
 * | Copyright (c) 2022 ivatin.com, Inc. All Rights Reserved
 * |-----------------------------------------------------------
 * | File: config.go
 * | Author: wuzhipeng at <wu.zhi.peng@outlook.com>
 * | Created: 2022-11-06
 * | Description: config.go
 * |-----------------------------------------------------------
 */

package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/ibytechaos/chaos/g"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Config is the configuration for the server.
type Config struct {
	// Server The address to listen on for HTTP requests.
	Server *Server
	// Client 客户端配置
	Client *Client
	// App 应用配置
	App *App
}

// Server 服务配置
type Server struct {
	// Log 日志配置
	Log *Log `mapstructure:"log" json:"log" yaml:"log"`
	// Hystrix 熔断配置
	Hystrix *Hystrix `mapstructure:"hystrix" json:"hystrix" yaml:"hystrix"`
	// Host 服务地址
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	// Port 端口
	Port int `mapstructure:"port" json:"port" yaml:"port"`
	// Domain 域名
	Domain string `mapstructure:"domain" json:"domain" yaml:"domain"`
}

// Address 获取服务地址
func (o *Server) Address() string {
	return o.Host + ":" + strconv.Itoa(o.Port)
}

// Hystrix 客户端熔断配置
type Hystrix struct {
	// Name 熔断器名称
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	// Timeout 超时时间
	Timeout int `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	// MaxConcurrent 最大并发数
	MaxConcurrent int `mapstructure:"max_concurrent" json:"max_concurrent" yaml:"max_concurrent"`
	// RequestVolumeThreshold 请求阈值
	RequestVolumeThreshold int `mapstructure:"request_volume_threshold" json:"request_volume_threshold" yaml:"request_volume_threshold"`
	// SleepWindow 熔断时间
	SleepWindow int `mapstructure:"sleep_window" json:"sleep_window" yaml:"sleep_window"`
	// ErrorPercentThreshold 错误百分比
	ErrorPercentThreshold int `mapstructure:"error_percent_threshold" json:"error_percent_threshold" yaml:"error_percent_threshold"`
}

// Log 日志配置
type Log struct {
	// Level 日志级别
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	// File 日志文件
	File string `mapstructure:"file" json:"file" yaml:"file"`
}

// Client 客户端配置
type Client struct {
	// MySQL MySQL配置
	MySQL *Channel `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	// Redis Redis配置
	Redis *Channel `mapstructure:"redis" json:"redis" yaml:"redis"`
	// Mongodb mongodb配置
	Mongodb *Channel `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
	// Smtp 邮件配置
	Smtp *Channel `mapstructure:"smtp" json:"smtp" yaml:"smtp"`
}

// Channel 通道配置
type Channel struct {
	// Hystrix 熔断配置
	Hystrix *Hystrix `mapstructure:"hystrix" json:"hystrix" yaml:"hystrix"`
	// Hosts 通道地址
	Hosts []string `mapstructure:"hosts" json:"hosts" yaml:"hosts"`
	// Host 通道地址
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	// Port 通道端口
	Port int `mapstructure:"port" json:"port" yaml:"port"`
	// Token 通道token
	Token string `mapstructure:"token" json:"token" yaml:"token"`
	// Name 通道名称
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	// Database 数据库名
	Database string `mapstructure:"database" json:"database" yaml:"database"`
	// Timeout 超时时间
	Timeout int `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	// MasterName 主节点
	MasterName string `mapstructure:"master_name" json:"master_name" yaml:"master_name"`
	// Username 用户名
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	// Password 密码
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}

// handle 补充配置
func (o *Channel) handle(name string) {
	if o == nil {
		return
	}
	o.Name = name
	if len(o.Hosts) > 0 && o.Host == "" {
		host := strings.Split(o.Hosts[0], ":")
		if len(host) == 2 {
			o.Host = host[0]
			o.Port, _ = strconv.Atoi(host[1])
		}
	}
	if o.Host != "" && o.Port > 0 && len(o.Hosts) == 0 {
		o.Hosts = []string{o.Host + ":" + strconv.Itoa(o.Port)}
	}
	if o.Token == "" {
		o.Token = o.Name
	}
	if o.Hystrix == nil {
		o.Hystrix = &Hystrix{
			Name:                   o.Name,
			Timeout:                o.Timeout,
			MaxConcurrent:          100,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  50,
		}
	}
}

// hystrix 补充配置
func (o *Server) hystrix() {
	if o.Hystrix == nil {
		o.Hystrix = &Hystrix{
			Name:                   "server",
			Timeout:                1000,
			MaxConcurrent:          100,
			RequestVolumeThreshold: 20,
			SleepWindow:            5000,
			ErrorPercentThreshold:  50,
		}
	}
}

// App 应用配置
type App struct {
	// Captcha 验证码配置
	Captcha *Captcha `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
}

// Captcha 验证码配置
type Captcha struct {
	// Height 高度
	Height int `mapstructure:"height" json:"height" yaml:"height"`
	// Width 宽度
	Width int `mapstructure:"width" json:"width" yaml:"width"`
	// Length 长度
	Length int `mapstructure:"length" json:"length" yaml:"length"`
}

var (
	instance *Config
	once     sync.Once
)

// GetInstance 获取配置实例
func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
		instance.init()
	})
	return instance
}

func (o *Config) init() {
	// 配置读取
	if err := NewViper("server").Watch(o.InitServer); err != nil {
		slog.Error("init server config fail", err)
	}
	if err := NewViper("client").Watch(o.InitClient); err != nil {
		slog.Error("init client config fail", err)
	}
	if err := NewViper("app").Watch(o.InitApp); err != nil {
		slog.Error("init app config fail", err)
	}
}

// InitServer 初始化服务配置
func (o *Config) InitServer(cfg *viper.Viper) error {
	if err := cfg.ReadInConfig(); err != nil {
		slog.Error("read config fail", err)
		return err
	}
	server := &Server{}
	if err := cfg.Unmarshal(server); err != nil {
		slog.Error("unmarshal config fail", err)
		return err
	}
	if o.Server == nil || o.Server.Log.Level != server.Log.Level {
		level := slog.LevelInfo
		if server.Log.Level == "debug" {
			level = slog.LevelDebug
		}
		// 初始化日志配置
		opts := slog.HandlerOptions{
			AddSource: true,
			Level:     level,
		}
		slog.SetDefault(slog.New(opts.NewTextHandler(os.Stderr)))
		slog.Default().Handler()
	}
	server.hystrix()
	o.Server = server
	return nil
}

// InitClient 初始化客户端配置
func (o *Config) InitClient(cfg *viper.Viper) error {
	if err := cfg.ReadInConfig(); err != nil {
		slog.Error("read config fail", err)
	}
	client := &Client{}
	if err := cfg.Unmarshal(client); err != nil {
		slog.Error("unmarshal config fail", err)
	}
	t := reflect.TypeOf(*client)
	v := reflect.ValueOf(*client)
	for i := 0; i < t.NumField(); i++ {
		v.Field(i).Interface().(*Channel).handle(t.Field(i).Name)
	}
	o.Client = client
	return nil
}

// InitApp 初始化应用配置
func (o *Config) InitApp(cfg *viper.Viper) error {
	if err := cfg.ReadInConfig(); err != nil {
		slog.Error("read config fail", err)
	}
	app := &App{}
	if err := cfg.Unmarshal(app); err != nil {
		slog.Error("unmarshal config fail", err)
	}
	o.App = app
	return nil
}

// Viper viper配置
type Viper struct {
	// Path 配置文件路径
	Path string
	// Name 配置文件名称
	Name string
	// Type 配置文件类型
	Type   string
	update time.Time
}

func NewViper(name string) *Viper {
	return &Viper{
		Path:   g.ConfDir,
		Name:   name,
		Type:   "yaml",
		update: time.Now(),
	}
}

// Watch 监听配置文件变化
func (o *Viper) Watch(f func(cfg *viper.Viper) error) error {
	v := viper.New()
	v.SetConfigName(o.Name)
	v.AddConfigPath(o.Path)
	v.SetConfigType(o.Type)
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		if o.update.Add(time.Second).After(time.Now()) {
			return
		}
		if err := f(v); err != nil {
			slog.Error("config update fail", err)
		}
		o.update = time.Now()
	})
	err := f(v)
	o.update = time.Now()
	return err
}
