package config

import (
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gogf/gf/v2/net/ghttp"
)

type Config struct {
	Servers struct {
		Inbound *ServerConfig
		Bizapi  *ServerConfig
		Admin   *ServerConfig
	}
	Logger   *LoggerConfig
	Database struct {
		Mysql *MysqlDaoConfig
		Mongo *MongoDaoConfig
		Redis *RedisDaoConfig
	}
	Broker     *BrokerConfig
	EventQueue struct {
		PushTaskEventConsumer *EventConsumerConfig
		PushLogEventConsumer  *EventConsumerConfig
	}
	Misc struct {
		Qos              string
		AuthKeyLength    int
		AuthSecretLength int
		LogExpireTime    int
		TokenExpireTime  int
	}
}

type ServerConfig struct {
	Name    string
	Address string
}

func (c *ServerConfig) Convert() *ghttp.ServerConfig {
	config := ghttp.NewConfig()
	config.Address = c.Address
	config.Name = c.Name
	return &config
}

type LoggerConfig struct {
	Path string
	File string
}

type MysqlDaoConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string // database name
}

type MongoDaoConfig struct {
	User string
	Pass string
	Host string
	Port string
	Name string // database name
}

func (c *MongoDaoConfig) URI() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%s", c.User, c.Pass, c.Host, c.Port)
}

type RedisDaoConfig struct {
	Pass string
	Host string
	Port string
}

func (c *RedisDaoConfig) Options() *redis.Options {
	return &redis.Options{
		Addr:     c.Host + ":" + c.Port,
		Password: c.Pass,
	}
}

type BrokerConfig struct {
	Broker  string
	Port    string
	Timeout int64 // second
}

func (c *BrokerConfig) BrokerAddress() string {
	return "tcp://" + c.Broker + ":" + c.Port
}

func (c *BrokerConfig) WaitTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

type EventConsumerConfig struct {
	Topic    string
	Dispatch uint
}

func (c *EventConsumerConfig) Check() bool {
	return c.Topic != "" && c.Dispatch > 0
}
