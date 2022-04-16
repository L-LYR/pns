package config

import (
	"context"
	"time"
)

type BrokerConfig struct {
	Broker  string
	Port    string
	Timeout int64 // second
}

func (c *BrokerConfig) BrokerAddress() string {
	return "tcp://" + c.Broker + ":" + c.Port
}

func (c *BrokerConfig) WaitTimeout() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

func MustLoadMQTTBrokerConfig(ctx context.Context) *BrokerConfig {
	brokerConfig := &BrokerConfig{}
	MustLoadConfig(ctx, "module.pusher.mqtt", brokerConfig)
	return brokerConfig
}
