package mqtt

import (
	"context"
	"errors"
	"fmt"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/util"
	"github.com/L-LYR/pns/proto/pkg/message"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/proto"
)

type Qos = byte

const (
	AtMostOnce  Qos = 0
	AtLeastOnce Qos = 1
	ExactlyOnce Qos = 2
)

type Client struct {
	Name         string
	Key          string
	Secret       string
	BrokerConfig *config.BrokerConfig
	Options      *paho.ClientOptions

	Client paho.Client
}

func MustNewClient(
	ctx context.Context,
	name string,
	key string,
	secret string,
	brokerConfig *config.BrokerConfig,
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
	options.AddBroker(p.BrokerConfig.BrokerAddress())
	options.SetClientID(p.Name)
	options.SetUsername(p.Key)
	options.SetPassword(p.Secret)
	options.SetConnectTimeout(p.BrokerConfig.WaitTimeout())
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
	pusherConfig *model.MQTTConfig,
	brokerConfig *config.BrokerConfig,
) *Client {
	return MustNewClient(
		ctx,
		util.GeneratePusherClientID(appId),
		pusherConfig.PusherKey,
		pusherConfig.PusherSecret,
		brokerConfig,
	)
}

func (p *Client) TryConnect() error {
	if !p.Client.IsConnected() {
		if token := p.Client.Connect(); token.WaitTimeout(p.Options.ConnectTimeout) {
			return nil
		} else if err := token.Error(); err != nil {
			return err
		} else {
			return errors.New("connection timeout")
		}
	}
	return nil
}

func (p *Client) Handle(ctx context.Context, task *model.PushTask) error {
	if err := p.TryConnect(); err != nil {
		return err
	}
	topic := _TopicOf(task)
	util.GLog.Infof(ctx, "topic: %s", topic)

	message := &message.Message{}
	if err := copier.Copy(message, task.Message); err != nil {
		return err
	}
	payload, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	token := p.Client.Publish(topic, AtMostOnce, false, payload)
	if ok := token.WaitTimeout(p.BrokerConfig.WaitTimeout()); !ok {
		return errors.New("message publish timeout")
	}
	return token.Error()
}

func _TopicOf(task *model.PushTask) string {
	switch task.Type {
	case model.PersonalPush:
		return fmt.Sprintf("%s/%d/%s/%d", task.Type.Name(), task.App.ID, task.Device.ID, task.ID)
	case model.BroadcastPush:
		return fmt.Sprintf("%s/%d/%d", task.Type.Name(), task.App.ID, task.ID)
	default:
		panic("unreachable")
	}
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
