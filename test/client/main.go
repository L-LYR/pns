package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/mobile/push_sdk/net/http"
	"github.com/L-LYR/pns/mobile/push_sdk/util"
	"github.com/L-LYR/pns/proto/pkg/message"
	paho "github.com/eclipse/paho.mqtt.golang"
	"google.golang.org/protobuf/proto"
)

// preset parameters
const (
	MaxDeviceId    = 1000000
	TestDuration   = 60 * time.Minute
	ConnectTimeout = 60 * time.Second
)

var (
	GlobalClient *http.Client
)

func initialize(ctx context.Context) {
	GlobalClient = http.MustNewHTTPClient("http://192.168.1.2:10086")
}

func ReportLog(payload http.Payload) {
	_, err := GlobalClient.POST("/log", payload)
	if err != nil {
		log.Printf("Error: %s", err.Error())
	}
}

type MessageHandler func(*message.Message) error

type Client struct {
	c paho.Client
}

func DefaultOptions(deviceId int) *paho.ClientOptions {
	o := paho.NewClientOptions()
	o.AddBroker("tcp://192.168.1.2:1883")
	o.SetClientID(util.GenerateClientId("pns-target", strconv.FormatInt(int64(deviceId), 10), 12345))
	o.SetUsername("test_app_name")
	o.SetPassword("test_app_name")
	o.SetConnectTimeout(ConnectTimeout)
	return o
}

func MustNewMQTTClient(deviceId int) *Client {
	p := &Client{c: paho.NewClient(DefaultOptions(deviceId))}
	if err := p.TryConnect(); err != nil {
		log.Printf("Error: %s", err.Error())
	}
	return p
}

func (c *Client) TryConnect() error {
	if !c.c.IsConnected() {
		token := c.c.Connect()
		if ok := token.WaitTimeout(ConnectTimeout); !ok {
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

func (c *Client) newEventLog(topic string, where string, err error) map[string]interface{} {
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
	case "DPush": // ss[1]: app id, ss[2]: device id, ss[3]: task id
		eventLog["appId"] = ss[1]
		eventLog["deviceId"] = ss[2]
		eventLog["taskId"] = ss[3]
	default:
		panic("unreachable")
	}
	return eventLog
}

func (c *Client) wrapHandler(fn MessageHandler) paho.MessageHandler {
	return func(_ paho.Client, m paho.Message) {
		msg, err := unmarshal(m)
		ReportLog(c.newEventLog(m.Topic(), "receive", err))
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
		err = fn(msg)
		ReportLog(c.newEventLog(m.Topic(), "show", err))
		if err != nil {
			log.Printf("Error: %s", err.Error())
			return
		}
	}
}

func (c *Client) subscribe(topic string, fn MessageHandler) {
	if err := c.TryConnect(); err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	token := c.c.Subscribe(topic, model.AtMostOnce, c.wrapHandler(fn))
	if ok := token.WaitTimeout(ConnectTimeout); !ok {
		log.Printf("Error: subscribe timeout")
		return
	} else if err := token.Error(); err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	log.Printf("Success to subscribe %s", topic)
}

func Topic(deviceId int) string {
	return fmt.Sprintf("DPush/12345/%d/+", deviceId)
}

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	clients := make([]*Client, 0, MaxDeviceId)
	initialize(ctx)
	for id := 0; id < 10000; id++ {
		c := MustNewMQTTClient(id)
		c.subscribe(Topic(id), func(m *message.Message) error {
			log.Printf("receive: %s:%s", m.Title, m.Content)
			return nil
		})
		clients = append(clients, c)
	}

	<-signalChan
	cancel()
	wg.Wait()
}
