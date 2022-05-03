package mqtt

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/target_status"
	"github.com/L-LYR/pns/internal/util"
	"github.com/L-LYR/pns/proto/pkg/message"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	Name         string
	Key          string
	Secret       string
	BrokerConfig *config.BrokerConfig
	Options      *paho.ClientOptions

	Client paho.Client
}

func _MustNewClient(
	ctx context.Context,
	appId int,
	name string,
	key string,
	secret string,
	brokerConfig *config.BrokerConfig,
) *Client {
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
	options.SetOnConnectHandler(_OnConnect(ctx, p.Name, appId))
	options.SetConnectionLostHandler(_OnConnectLost(ctx, p.Name))
	options.SetReconnectingHandler(_OnReconnecting(ctx))
	p.Options = options
	p.Client = paho.NewClient(options)
	onlineTopic := fmt.Sprintf("PNS/online/%d/+", appId)
	offlineTopic := fmt.Sprintf("PNS/offline/%d/+", appId)
	if err := p.Subscribe(ctx, onlineTopic); err != nil {
		util.GLog.Warningf(ctx, "Client %s fail to subscribe %s, because %s", p.Name, onlineTopic, err.Error())
	}
	if err := p.Subscribe(ctx, offlineTopic); err != nil {
		util.GLog.Warningf(ctx, "Client %s fail to subscribe %s, because %s", p.Name, offlineTopic, err.Error())
	}
	return p
}

func MustNewPusher(
	ctx context.Context,
	appId int,
	pusherConfig *model.MQTTConfig,
	brokerConfig *config.BrokerConfig,
) *Client {
	return _MustNewClient(
		ctx,
		appId,
		util.GeneratePusherClientID(appId),
		pusherConfig.PusherKey,
		pusherConfig.PusherSecret,
		brokerConfig,
	)
}

func (p *Client) TryConnect() error {
	if !p.Client.IsConnected() {
		if token := p.Client.Connect(); !token.WaitTimeout(p.Options.ConnectTimeout) {
			return errors.New("connection timeout")
		} else if err := token.Error(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Client) Handle(ctx context.Context, task model.PushTask) error {
	if err := p.TryConnect(); err != nil {
		return err
	}
	topic := task.GetTopic()
	util.GLog.Infof(ctx, "topic: %s", topic)

	message := &message.Message{
		ID: int64(task.GetID()),
	}
	if err := copier.Copy(message, task.GetMessage()); err != nil {
		return err
	}
	payload, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	token := p.Client.Publish(topic, task.GetQos(), false, payload)
	if ok := token.WaitTimeout(p.BrokerConfig.WaitTimeout()); !ok {
		return errors.New("message publish timeout")
	}
	return token.Error()
}

func (p *Client) Close(context.Context) {
	p.Client.Disconnect(100)
}

func (p *Client) Subscribe(ctx context.Context, topic string) error {
	if err := p.TryConnect(); err != nil {
		return err
	}

	token := p.Client.Subscribe(topic, model.AtMostOnce,
		func(c paho.Client, m paho.Message) {
			p.targetStatusTrace(ctx, m.Topic())
		},
	)
	if ok := token.WaitTimeout(p.Options.ConnectTimeout); !ok {
		return errors.New("subscribe timeout")
	} else if err := token.Error(); err != nil {
		return err
	}
	util.GLog.Printf(ctx, "Subscribe %s successfully", topic)
	return nil
}

func (p *Client) targetStatusTrace(ctx context.Context, topic string) {
	ss := strings.Split(topic, "/")
	if len(ss) != 4 {
		util.GLog.Warningf(ctx, "Cannot parse %s as status trace", topic)
		return
	}
	var err error
	switch ss[1] {
	case "online":
		err = target_status.SetTargetOnline(ctx, ss[2], ss[3])
	case "offline":
		err = target_status.SetTargetOffline(ctx, ss[2], ss[3])
	default:
		util.GLog.Warningf(ctx, "Unknown status: %s", ss[1])
	}
	if err != nil {
		util.GLog.Errorf(ctx, "Fail to set status for %s, because %s", topic, err)
	}
}

func _OnConnect(ctx context.Context, name string, appId int) paho.OnConnectHandler {

	return func(client paho.Client) {
		util.GLog.Noticef(ctx, "Client %s connect with broker", name)
	}
}

func _OnConnectLost(ctx context.Context, name string) paho.ConnectionLostHandler {
	return func(client paho.Client, err error) {
		util.GLog.Noticef(ctx, "Client %s lost connection with broker, because:%+v", name, err)
	}
}

func _OnReconnecting(ctx context.Context) paho.ReconnectHandler {
	return func(client paho.Client, options *paho.ClientOptions) {
		util.GLog.Noticef(ctx, "Client %s is connecting with broker", options.ClientID)
	}
}
