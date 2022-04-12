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
	PersonalTopic  string
	BroadcastTopic string

	Handlers map[string]paho.MessageHandler
}

func NewTopicSet(cfg *storage.Config) *TopicSet {
	ts := &TopicSet{Handlers: make(map[string]paho.MessageHandler)}
	if cfg != nil {
		ts.PersonalTopic = fmt.Sprintf("PPush/%d/%s/+", cfg.AppId, cfg.DeviceId)
		ts.BroadcastTopic = fmt.Sprintf("BPush/%d/+", cfg.AppId)
	}
	return ts
}

type LogHandler func(fmt string, v ...interface{})

type MessageHandler func(*message.Message) error

func _NewEventLog(topic string, where string, err error) map[string]interface{} {
	eventLog := make(map[string]interface{})
	if err != nil {
		eventLog["hint"] = fmt.Sprintf("failure: %s", err)
	} else {
		eventLog["hint"] = "success"
	}
	eventLog["where"] = where
	eventLog["timestamp"] = time.Now().UnixMilli()
	ss := strings.Split(topic, "/")
	switch ss[0] {
	case "PPush": // ss[1]: app id, ss[2]: device id, ss[3]: task id
		eventLog["appId"] = ss[1]
		eventLog["deviceId"] = ss[2]
		eventLog["taskId"] = ss[3]
	default:
		panic("Unreachable")
	}
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

func (o *Options) SetWithCfg(cfg *storage.Config) {
	o.AddBroker(cfg.GetAddress())
	o.SetClientID(cfg.ClientId)
	o.SetUsername(cfg.Key)
	o.SetPassword(cfg.Secret)
	o.SetConnectTimeout(cfg.GetConnectTimeout())
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
			// NOTICE: in first connection, the handlers is empty
			for topic, handler := range options.topicSet.Handlers {
				p.subscribe(topic, handler)
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

	return p
}

func (c *Client) SetLogHandler(fn LogHandler) {
	c.options.SetLogHandler(fn)
}

func (c *Client) TryConnect() error {
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

func unmarshal(m paho.Message) (*message.Message, error) {
	message := &message.Message{}
	if err := proto.Unmarshal(m.Payload(), message); err != nil {
		return nil, err
	}
	return message, nil
}

func (c *Client) _WrapHandler(fn MessageHandler) paho.MessageHandler {
	return func(_ paho.Client, m paho.Message) {
		msg, err := unmarshal(m)
		defer c.options.recvHandler(_NewEventLog(m.Topic(), "receive", err))
		if err != nil {
			c.options.logHandler("Error: %s", err.Error())
			return
		}
		err = fn(msg)
		defer c.options.showHandler(_NewEventLog(m.Topic(), "show", err))
		if err != nil {
			c.options.logHandler("Error: %s", err.Error())
			return
		}
	}
}

func (c *Client) SubscribePersonalPush(fn MessageHandler) {
	c.subscribe(c.options.topicSet.PersonalTopic, c._WrapHandler(fn))
}

func (c *Client) SubscribeBroadcastPush(fn MessageHandler) {
	c.subscribe(c.options.topicSet.BroadcastTopic, c._WrapHandler(fn))
}

func (c *Client) subscribe(topic string, fn paho.MessageHandler) {
	if err := c.TryConnect(); err != nil {
		c.options.logHandler("Error: %s", err.Error())
		return
	}
	c.options.topicSet.Handlers[topic] = fn
	if token := c.c.Subscribe(topic, AtMostOnce, fn); !token.WaitTimeout(c.options.ConnectTimeout) {
		c.options.logHandler("Error: subscribe timeout")
		return
	} else if err := token.Error(); err != nil {
		c.options.logHandler("Error: %s", err.Error())
		return
	}
}
