package config

import (
	"context"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var (
	_Config *Config = &Config{}
)

func MustLoad(ctx context.Context) {
	if !g.Cfg().Available(ctx) {
		panic("global config is not available")
	}

	if v, err := g.Cfg().Get(ctx, "."); err != nil {
		util.GLog.Warningf(ctx, "Fail to load config, use default")
	} else if err := v.Struct(_Config); err != nil {
		util.GLog.Warningf(ctx, "Fail to load config, use default")
	} else {
		util.GLog.Info(ctx, "Success to load config")
		return
	}

	_Config = DefaultConfig()
}

func CommonTaskQos() model.Qos {
	return model.ParseQos(_Config.Misc.Qos)
}

func PushTaskEventConsumerConfig() *EventConsumerConfig {
	return _Config.EventQueue.PushTaskEventConsumer
}

func PushLogEventConsumerConfig() *EventConsumerConfig {
	return _Config.EventQueue.PushLogEventConsumer
}

func InboundServerConfig() *ghttp.ServerConfig {
	return _Config.Servers.Inbound.Convert()
}
func BizapiServerConfig() *ghttp.ServerConfig {
	return _Config.Servers.Bizapi.Convert()
}
func AdminServerConfig() *ghttp.ServerConfig {
	return _Config.Servers.Admin.Convert()
}

func PushTaskEventTopic() string {
	return _Config.EventQueue.PushTaskEventConsumer.Topic
}
func PushLogEventTopic() string {
	return _Config.EventQueue.PushLogEventConsumer.Topic
}

func MQTTBrokerConfig() *BrokerConfig {
	return _Config.Broker
}

func RedisConfig() *RedisDaoConfig {
	return _Config.Database.Redis
}

func MongoConfig() *MongoDaoConfig {
	return _Config.Database.Mongo
}

func MysqlConfig() *MysqlDaoConfig {
	return _Config.Database.Mysql
}

func AuthKeyLength() int {
	return _Config.Misc.AuthKeyLength
}

func AuthSecretLength() int {
	return _Config.Misc.AuthSecretLength
}

func LogExpireTime() time.Duration {
	return time.Second * time.Duration(_Config.Misc.LogExpireTime)
}

func TokenExpireTime() time.Duration {
	return time.Second * time.Duration(_Config.Misc.TokenExpireTime)
}
