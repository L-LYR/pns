package mqtt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/L-LYR/pns/mobile/push_sdk/storage"
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
	AppId          int
	DeviceId       string
	DirectTopic    string
	BroadcastTopic string

	Handlers map[string]paho.MessageHandler
}

func NewTopicSet(cfg *storage.Config) *TopicSet {
	if cfg == nil {
		panic("config is nil")
	}
	ts := &TopicSet{
		AppId:    cfg.App.ID,
		DeviceId: cfg.DeviceId,
		Handlers: make(map[string]paho.MessageHandler),
	}
	ts.DirectTopic = fmt.Sprintf("DPush/%d/%s", ts.AppId, ts.DeviceId)
	ts.BroadcastTopic = fmt.Sprintf("BPush/%d", ts.AppId)
	return ts
}

type LogHandler func(fmt string, v ...interface{})

type MessageHandler func(*message.Message) error

func (c *Client) newEventLog(topic string, where string, message *message.Message, err error) map[string]interface{} {
	eventLog := make(map[string]interface{})
	if err != nil {
		eventLog["hint"] = fmt.Sprintf("failed, because %s", err.Error())
	} else {
		eventLog["hint"] = "success"
	}
	eventLog["where"] = where
	eventLog["timestamp"] = time.Now().UnixMilli()
	ss := strings.Split(topic, "/")
	switch ss[0] {
	case "DPush": // ss[1]: app id, ss[2]: device id
		eventLog["appId"] = ss[1]
		eventLog["deviceId"] = ss[2]
	case "BPush":
		eventLog["appId"] = ss[1]
		eventLog["deviceId"] = c.options.topicSet.DeviceId
	default:
		panic("unreachable")
	}
	eventLog["taskId"] = message.ID
	return eventLog
}

type EventLogHandler func(map[string]interface{})

type Options struct {
	*paho.ClientOptions

	topicSet *TopicSet

	logHandler  LogHandler
	recvHandler EventLogHandler
	showHandler EventLogHandler
}

func NewOptions() *Options {
	return &Options{
		ClientOptions: paho.NewClientOptions(),
	}
}

func _OfflineTopic(appId int, deviceId string) string {
	return fmt.Sprintf("PNS/offline/%d/%s", appId, deviceId)
}

func _OnlineTopic(appId int, deviceId string) string {
	return fmt.Sprintf("PNS/online/%d/%s", appId, deviceId)
}

func (o *Options) SetWithCfg(cfg *storage.Config) {
	o.AddBroker(cfg.GetAddress())
	o.SetClientID(cfg.ClientId)
	o.SetUsername(cfg.App.Key)
	o.SetPassword(cfg.App.Secret)
	o.SetConnectTimeout(cfg.GetConnectTimeout())
	o.SetWill(_OfflineTopic(cfg.GetAppId(), cfg.GetDeviceId()), "", ExactlyOnce, false)
	o.SetCleanSession(false)
}

func (o *Options) SetLogHandler(fn LogHandler) {
	o.logHandler = fn
}

func (o *Options) SetTopicSet(ts *TopicSet) {
	o.topicSet = ts
}

func (o *Options) SetRecvHandler(fn EventLogHandler) {
	o.recvHandler = fn
}

func (o *Options) SetShowHandler(fn EventLogHandler) {
	o.showHandler = fn
}

type Client struct {
	options Options
	c       paho.Client
}

func MustNewMQTTClient(options *Options) *Client {
	p := &Client{
		options: *options,
	}

	options.ClientOptions.SetOnConnectHandler(
		func(c paho.Client) {
			options.logHandler("Info: connected")
			for topic, handler := range options.topicSet.Handlers {
				p.subscribe(topic, handler)
			}
			// publish online message
			token := c.Publish(
				_OnlineTopic(
					options.topicSet.AppId,
					options.topicSet.DeviceId,
				),
				ExactlyOnce, false, nil,
			)
			if ok := token.WaitTimeout(options.ConnectTimeout); !ok {
				options.logHandler("Error: message publish timeout")
			}
		},
	)

	options.ClientOptions.SetConnectionLostHandler(
		func(c paho.Client, err error) {
			options.logHandler("Info: connection lost")
		},
	)

	options.ClientOptions.SetReconnectingHandler(
		func(c paho.Client, o *paho.ClientOptions) {
			options.logHandler("Info: reconnecting")
		},
	)

	p.c = paho.NewClient(options.ClientOptions)

	if err := p.TryConnect(); err != nil {
		p.options.logHandler("Error: %s", err.Error())
	}

	return p
}

func (c *Client) SetLogHandler(fn LogHandler) {
	c.options.SetLogHandler(fn)
}

func (c *Client) TryConnect() error {
	if !c.c.IsConnected() {
		token := c.c.Connect()
		if !token.WaitTimeout(c.options.ConnectTimeout) {
			return errors.New("connection timeout")
		} else if err := token.Error(); err != nil {
			return err
		}
	}
	return nil
}

func unmarshal(m paho.Message) (*message.Message, error) {
	message := &message.Message{}
	if err := proto.Unmarshal(m.Payload(), message); err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Client) wrapHandler(fn MessageHandler) paho.MessageHandler {
	return func(_ paho.Client, m paho.Message) {
		msg, err := unmarshal(m)
		defer c.options.recvHandler(c.newEventLog(m.Topic(), "receive", msg, err))
		if err != nil {
			c.options.logHandler("Error: %s", err.Error())
			return
		}
		err = fn(msg)
		defer c.options.showHandler(c.newEventLog(m.Topic(), "show", msg, err))
		if err != nil {
			c.options.logHandler("Error: %s", err.Error())
			return
		}
	}
}

func (c *Client) SubscribePersonalPush(fn MessageHandler) {
	c.subscribe(c.options.topicSet.DirectTopic, c.wrapHandler(fn))
}

func (c *Client) SubscribeBroadcastPush(fn MessageHandler) {
	c.subscribe(c.options.topicSet.BroadcastTopic, c.wrapHandler(fn))
}

func (c *Client) subscribe(topic string, fn paho.MessageHandler) {
	if err := c.TryConnect(); err != nil {
		c.options.logHandler("Error: %s", err.Error())
		return
	}
	c.options.topicSet.Handlers[topic] = fn
	token := c.c.Subscribe(topic, AtMostOnce, fn)
	if !token.WaitTimeout(c.options.ConnectTimeout) {
		c.options.logHandler("Error: subscribe timeout")
		return
	} else if err := token.Error(); err != nil {
		c.options.logHandler("Error: %s", err.Error())
		return
	}
	c.options.logHandler("Subscribe %s successfully", topic)
}
