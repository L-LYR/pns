package net

import (
	"context"
	"errors"
	"fmt"

	"github.com/L-LYR/pns/mobile/push_sdk/storage"
	"github.com/L-LYR/pns/mobile/push_sdk/util"
	"github.com/L-LYR/pns/proto/pkg/message"
	paho "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
)

type Qos = byte

const (
	AtMostOnce  Qos = 0
	AtLeastOnce Qos = 1
	ExactlyOnce Qos = 2
)

type TopicSet struct {
	PersonalTopic  string
	BroadcastTopic string
}

func (ts *TopicSet) Register(cfg *storage.Config) {
	ts.PersonalTopic = fmt.Sprintf("PPush/%d/%s/+", cfg.AppId, cfg.DeviceId)
	ts.BroadcastTopic = fmt.Sprintf("BPush/%d/+", cfg.AppId)
}

type MQTTClient struct {
	topicSet *TopicSet
	options  *paho.ClientOptions
	c        paho.Client
}

func MustNewMQTTClient(ctx context.Context, cfg *storage.Config) *MQTTClient {
	options := paho.NewClientOptions()
	options.AddBroker(cfg.GetAddress())
	options.SetClientID(cfg.ClientID)
	options.SetUsername(cfg.Key)
	options.SetPassword(cfg.Secret)
	options.SetConnectTimeout(cfg.GetConnectTimeout())
	options.SetOnConnectHandler(_OnConnect(ctx))
	options.SetConnectionLostHandler(_OnConnectLost(ctx))
	options.SetReconnectingHandler(_OnReconnecting(ctx))
	topicSet := &TopicSet{}
	topicSet.Register(cfg)

	p := &MQTTClient{
		options:  options,
		c:        paho.NewClient(options),
		topicSet: topicSet,
	}
	if err := p.TryConnect(); err != nil {
		util.Log("Error: %s", err.Error())
	}
	return p
}

func (c *MQTTClient) TryConnect() error {
	if !c.c.IsConnected() {
		if token := c.c.Connect(); token.WaitTimeout(c.options.ConnectTimeout) {
			return nil
		} else if err := token.Error(); err != nil {
			return err
		} else {
			return errors.New("connection timeout")
		}
	}
	return nil
}

type MessageHandler func(*message.Message)

func unmarshal(m paho.Message) *message.Message {
	message := &message.Message{}
	if err := proto.Unmarshal(m.Payload(), message); err != nil {
		// TODO: Dump Local Error Log
		return nil
	}
	return message
}

func (c *MQTTClient) SubscribePersonalPush(fn MessageHandler) {
	handler := func(c paho.Client, m paho.Message) {
		fn(unmarshal(m))
	}
	c.subscribe(c.topicSet.PersonalTopic, handler)
}

func (c *MQTTClient) SubscribeBroadcastPushHandler(fn MessageHandler) {
	handler := func(c paho.Client, m paho.Message) {
		fn(unmarshal(m))
	}
	c.subscribe(c.topicSet.BroadcastTopic, handler)
}

func (c *MQTTClient) subscribe(topic string, fn paho.MessageHandler) {
	if err := c.TryConnect(); err != nil {
		util.Log("Error: %s", err.Error())
		return
	}
	if token := c.c.Subscribe(topic, AtMostOnce, fn); !token.WaitTimeout(c.options.ConnectTimeout) {
		util.Log("subscribe timeout")
		return
	} else if err := token.Error(); err != nil {
		util.Log("Error: %s", err.Error())
		return
	}
	util.Log("Info: Subscribe %s topic successfully", topic)
}

func _OnConnect(ctx context.Context) paho.OnConnectHandler {
	return func(client paho.Client) {
		util.Log("Notice: mqtt connected")
	}
}

func _OnConnectLost(ctx context.Context) paho.ConnectionLostHandler {
	return func(client paho.Client, err error) {
		util.Log("Notice: mqtt lost connection")
	}
}

func _OnReconnecting(ctx context.Context) paho.ReconnectHandler {
	return func(client paho.Client, options *paho.ClientOptions) {
		util.Log("Notice: mqtt reconnecting")
	}
}
