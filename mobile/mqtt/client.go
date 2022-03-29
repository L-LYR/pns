package mqtt

import (
	"errors"
	"fmt"
	"time"

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

type ClientConfig struct {
	Name string

	AppId  int
	Key    string
	Secret string

	DeviceId string

	Broker  string
	Port    string
	Timeout int64
}

func (c *ClientConfig) address() string {
	return "tcp://" + c.Broker + ":" + c.Port
}

func (c *ClientConfig) timeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

type Client struct {
	Config  *ClientConfig
	Options *paho.ClientOptions
	Client  paho.Client
}

func MustNewClient(cfg *ClientConfig) *Client {
	p := &Client{Config: cfg}
	options := paho.NewClientOptions()
	options.AddBroker(p.Config.address())
	options.SetClientID(p.Config.Name)
	options.SetUsername(p.Config.Key)
	options.SetPassword(p.Config.Secret)
	// options.SetOnConnectHandler(_OnConnect(ctx, p.Config.Name))
	// options.SetConnectionLostHandler(_OnConnectLost(ctx, p.Config.Name))
	// options.SetReconnectingHandler(_OnReconnecting(ctx))
	p.Options = options
	p.Client = paho.NewClient(options)
	return p
}

func (c *Client) Topic() string {
	return fmt.Sprintf("Push/%d/%s/#", c.Config.AppId, c.Config.DeviceId)
}

type MessageHandler func(*message.Message)

func (c *Client) Subscribe(fn MessageHandler) error {
	if !c.Client.IsConnected() {
		token := c.Client.Connect()
		if ok := token.WaitTimeout(c.Config.timeout()); !ok {
			return errors.New("connect timeout")
		}
		if err := token.Error(); err != nil {
			return err
		}
	}

	token := c.Client.Subscribe(c.Topic(), AtMostOnce, func(c paho.Client, m paho.Message) {
		message := &message.Message{}
		if err := proto.Unmarshal(m.Payload(), message); err != nil {
			panic(err)
		}
		fn(message)
	})
	return token.Error()
}

// func _OnConnect(ctx context.Context, name string) paho.OnConnectHandler {
// 	return func(client paho.Client) {
// 	}
// }

// func _OnConnectLost(ctx context.Context, name string) paho.ConnectionLostHandler {
// 	return func(client paho.Client, err error) {
// 	}
// }

// func _OnReconnecting(ctx context.Context) paho.ReconnectHandler {
// 	return func(client paho.Client, options *paho.ClientOptions) {
// 	}
// }
