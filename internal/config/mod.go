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
		util.GLog.Panicf(ctx, "Global config is not available")
	}

	if v, err := g.Cfg().Get(ctx, "."); err != nil {
		util.GLog.Warningf(ctx, "Fail to load config, use default")
	} else if err := v.Struct(_Config); err != nil {
		util.GLog.Warningf(ctx, "Fail to load config, use default")
	} else {
		util.GLog.Infof(ctx, "Success to load config:\n%s", util.Readable(_Config))
		return
	}

	_Config = DefaultConfig(ctx)
}

func CommonTaskQos() model.Qos {
	return model.ParseQos(_Config.Misc.Qos)
}

func TaskValidationEventConsumerConfig() *EventQueueConfig {
	return _Config.EventQueue.TaskValidationEventQueue
}

func DirectPushTaskEventConsumerConfig() *EventQueueConfig {
	return _Config.EventQueue.DirectPushTaskEventQueue
}

func BroadcastPushTaskEventConsumerConfig() *EventQueueConfig {
	return _Config.EventQueue.BroadcastPushTaskEventQueue
}

func PushLogEventConsumerConfig() *EventQueueConfig {
	return _Config.EventQueue.PushLogEventqueue
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

func TaskValidationEventTopic() string {
	return _Config.EventQueue.TaskValidationEventQueue.Topic
}

func DirectPushTaskEventTopic() string {
	return _Config.EventQueue.DirectPushTaskEventQueue.Topic
}

func BroadcastPushTaskEventTopic() string {
	return _Config.EventQueue.BroadcastPushTaskEventQueue.Topic
}

func PushLogEventTopic() string {
	return _Config.EventQueue.PushLogEventqueue.Topic
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

func GetEnginePoolMinLen() int64 {
	return _Config.EnginePool.MinLen
}

func GetEnginePoolMaxLen() int64 {
	return _Config.EnginePool.MaxLen
}

func GetAppLevelFreqCtrlConfig() *FreqCtrlConfig {
	return _Config.FrequencyControl.AppLevel
}

func GetTargetLevelFreqCtrlConfig() *FreqCtrlConfig {
	return _Config.FrequencyControl.TargetLevel
}
