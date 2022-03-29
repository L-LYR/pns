package mqtt

import (
	"time"
)

type Qos = byte

const (
	AtMostOnce  Qos = 0
	AtLeastOnce Qos = 1
	ExactlyOnce Qos = 2
)

type BrokerConfig struct {
	Broker  string
	Port    string
	Timeout int64 // second
}

func (c *BrokerConfig) brokerAddress() string {
	return "tcp://" + c.Broker + ":" + c.Port
}

func (c *BrokerConfig) timeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}
