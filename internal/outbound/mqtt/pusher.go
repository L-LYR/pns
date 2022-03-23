package mqtt

import (
	"context"
	"errors"
	"time"

	"github.com/L-LYR/pns/internal/dynamic_config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/util"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	jsoniter "github.com/json-iterator/go"
)

const (
	_PusherConfigName = "module.pusher.mqtt"
)

type Qos = byte

const (
	AtMostOnce  Qos = 0
	AtLeastOnce Qos = 1
	ExactlyOnce Qos = 2
)

var (
	_DefaultPusherConfig *PusherConfig = nil
)

type PusherConfig struct {
	Broker  string
	Port    string
	Timeout int64 // second
}

func (pc *PusherConfig) ServerName() string {
	return "tcp://" + pc.Broker + ":" + pc.Port
}

func (pc *PusherConfig) WaitTimeout() time.Duration {
	return time.Duration(pc.Timeout) * time.Second
}

type Pusher struct {
	PusherConfig *PusherConfig
	Auth         *model.MQTTConfig
	Name         string
	Client       paho.Client
	Options      *paho.ClientOptions
}

func MustNewPusher(
	ctx context.Context,
	appConfig *model.AppConfig,
	auth *model.MQTTConfig,
	pusherConfig *PusherConfig,
) *Pusher {
	if auth == nil || pusherConfig == nil {
		panic("pusher config or app config is not given")
	}
	p := &Pusher{
		Name:         util.GeneratePusherClientID(appConfig.ID),
		Auth:         auth,
		PusherConfig: pusherConfig,
	}
	options := paho.NewClientOptions()
	options.AddBroker(p.PusherConfig.ServerName())
	options.SetClientID(p.Name)
	options.SetUsername(auth.Key)
	options.SetPassword(auth.Secret)
	options.SetOnConnectHandler(OnConnect(ctx, p.Name))
	options.SetConnectionLostHandler(OnConnectLost(ctx, p.Name))
	options.SetReconnectingHandler(OnReconnecting(ctx))
	p.Options = options
	p.Client = paho.NewClient(options)
	return p
}

func MustNewDefaultPusher(ctx context.Context, appId int) *Pusher {
	if _DefaultPusherConfig == nil {
		pusherConfig := &PusherConfig{}
		config := g.Cfg().MustGet(ctx, _PusherConfigName).Map()
		if err := gconv.Struct(config, pusherConfig); err != nil {
			panic("fail to initialize default pusher config")
		}
		_DefaultPusherConfig = pusherConfig
	}
	return MustNewPusher(
		ctx,
		dynamic_config.GetAppConfigByAppId(appId),
		dynamic_config.GetPusherAuthByAppId(appId, model.MQTTPusher).(*model.MQTTConfig),
		_DefaultPusherConfig,
	)
}

func (p *Pusher) Handle(ctx context.Context, task *model.PushTask) error {
	if !p.Client.IsConnected() {
		token := p.Client.Connect()
		if ok := token.WaitTimeout(p.PusherConfig.WaitTimeout()); !ok {
			return errors.New("connect timeout")
		}
		if err := token.Error(); err != nil {
			return err
		}
	}

	topic := TopicOf(task)
	util.GLog.Infof(ctx, "topic: %s", topic)
	payload, err := jsoniter.MarshalToString(task.Message)
	if err != nil {
		return err
	}
	token := p.Client.Publish(topic, AtMostOnce, false, payload)
	if ok := token.WaitTimeout(p.PusherConfig.WaitTimeout()); !ok {
		return errors.New("message publish timeout")
	}
	return token.Error()
}

// We simply determine that the topic of single push uses the topic with format: "/<app_name>/<device_token>"
func TopicOf(task *model.PushTask) string {
	appName := dynamic_config.GetAppNameByAppId(task.App.ID)
	return "/" + appName + "/" + task.Device.ID
}

func OnConnect(ctx context.Context, name string) paho.OnConnectHandler {
	return func(client paho.Client) {
		util.GLog.Noticef(ctx, "Client %s connect with broker", name)
	}
}

func OnConnectLost(ctx context.Context, name string) paho.ConnectionLostHandler {
	return func(client paho.Client, err error) {
		util.GLog.Noticef(ctx, "Client %s lost connection with broker, because:%+v", err, name)
	}
}

func OnReconnecting(ctx context.Context) paho.ReconnectHandler {
	return func(client paho.Client, options *paho.ClientOptions) {
		util.GLog.Noticef(ctx, "Client %s is connecting with broker", options.ClientID)
	}
}
