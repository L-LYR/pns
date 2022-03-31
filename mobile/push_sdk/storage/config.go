package storage

import "time"

type Config struct {
	ClientID string

	AppId  int
	Key    string
	Secret string

	DeviceId string
	Token    map[string]string

	Broker         string
	Port           string
	RetryInterval  int64
	ConnectTimeout int64
}

func (c *Config) GetAddress() string {
	return "tcp://" + c.Broker + ":" + c.Port
}

func (c *Config) GetRetryInterval() time.Duration {
	return time.Duration(c.RetryInterval) * time.Millisecond
}

func (c *Config) GetConnectTimeout() time.Duration {
	return time.Duration(c.ConnectTimeout) * time.Second
}
