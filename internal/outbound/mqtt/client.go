package mqtt

import (
	"context"
	"errors"
	"fmt"

	"github.com/L-LYR/pns/internal/local_storage"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/util"
	paho "github.com/eclipse/paho.mqtt.golang"
	jsoniter "github.com/json-iterator/go"
)

type Client struct {
	Name         string
	Key          string
	Secret       string
	BrokerConfig *BrokerConfig
	Options      *paho.ClientOptions

	Client paho.Client
}

func MustNewClient(
	ctx context.Context,
	name string,
	key string,
	secret string,
	brokerConfig *BrokerConfig,
) *Client {
	if brokerConfig == nil {
		panic("pusher config or app config is not given")
	}
	p := &Client{
		Name:         name,
		Key:          key,
		Secret:       secret,
		BrokerConfig: brokerConfig,
	}
	options := paho.NewClientOptions()
	options.AddBroker(p.BrokerConfig.brokerAddress())
	options.SetClientID(p.Name)
	options.SetUsername(p.Key)
	options.SetPassword(p.Secret)
	options.SetOnConnectHandler(_OnConnect(ctx, p.Name))
	options.SetConnectionLostHandler(_OnConnectLost(ctx, p.Name))
	options.SetReconnectingHandler(_OnReconnecting(ctx))
	p.Options = options
	p.Client = paho.NewClient(options)
	return p
}

func MustNewPusher(
	ctx context.Context,
	appId int,
	brokerConfig *BrokerConfig,
) *Client {
	auth := local_storage.GetPusherAuthByAppId(appId, model.MQTTPusher).(*model.MQTTConfig)
	return MustNewClient(
		ctx,
		util.GeneratePusherClientID(appId),
		auth.PusherKey,
		auth.PusherSecret,
		brokerConfig,
	)
}

// func MustNewReceiver(
// 	ctx context.Context,
// 	deviceId string,
// 	appId int,
// 	key string,
// 	secret string,
// 	brokerConfig *BrokerConfig,
// ) *Client {
// 	return MustNewClient(
// 		ctx,
// 		util.GenerateTargetClientID(deviceId, appId),
// 		key,
// 		secret,
// 		brokerConfig,
// 	)
// }

func (p *Client) Handle(ctx context.Context, task *model.PushTask) error {
	if !p.Client.IsConnected() {
		token := p.Client.Connect()
		if ok := token.WaitTimeout(p.BrokerConfig.timeout()); !ok {
			return errors.New("connect timeout")
		}
		if err := token.Error(); err != nil {
			return err
		}
	}

	topic := _TopicOf(task)
	util.GLog.Infof(ctx, "topic: %s", topic)
	payload, err := jsoniter.MarshalToString(task.Message)
	if err != nil {
		return err
	}
	token := p.Client.Publish(topic, AtMostOnce, false, payload)
	if ok := token.WaitTimeout(p.BrokerConfig.timeout()); !ok {
		return errors.New("message publish timeout")
	}
	return token.Error()
}

// We simply determine that the topic of single push uses the topic with format: "/<app_name>/<device_token>"
func _TopicOf(task *model.PushTask) string {
	return fmt.Sprintf("Push/%d/%s/%d", task.App.ID, task.Device.ID, task.ID)
}

func _OnConnect(ctx context.Context, name string) paho.OnConnectHandler {
	return func(client paho.Client) {
		util.GLog.Noticef(ctx, "Client %s connect with broker", name)
	}
}

func _OnConnectLost(ctx context.Context, name string) paho.ConnectionLostHandler {
	return func(client paho.Client, err error) {
		util.GLog.Noticef(ctx, "Client %s lost connection with broker, because:%+v", err, name)
	}
}

func _OnReconnecting(ctx context.Context) paho.ReconnectHandler {
	return func(client paho.Client, options *paho.ClientOptions) {
		util.GLog.Noticef(ctx, "Client %s is connecting with broker", options.ClientID)
	}
}
