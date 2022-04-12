package config

import "context"

type ConsumerConfig struct {
	Topic    string
	Dispatch uint
}

func (c *ConsumerConfig) Check() bool {
	return c.Topic != "" && c.Dispatch > 0
}

func MustLoadConsumerConfig(ctx context.Context, name string) *ConsumerConfig {
	cfg := &ConsumerConfig{}
	MustLoadConfig(ctx, "module.event_queue."+name, cfg)
	_SetTopicName(name, cfg.Topic)
	return cfg
}

var (
	_PushEventTopic   string
	_LogEventTopic string
)

func PushEventTopic() string   { return _PushEventTopic }
func LogEventTopic() string { return _LogEventTopic }

func _SetTopicName(cfgName string, topic string) {
	switch cfgName {
	case "push_event_consumer":
		_PushEventTopic = topic
	case "log_event_consumer":
		_LogEventTopic = topic
	default:
		panic("unreachable")
	}
}
