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
	_TargetEventTopic string
)

func PushEventTopic() string   { return _PushEventTopic }
func TargetEventTopic() string { return _TargetEventTopic }

func _SetTopicName(cfgName string, topic string) {
	switch cfgName {
	case "push_event_consumer":
		_PushEventTopic = topic
	case "target_event_consumer":
		_TargetEventTopic = topic
	default:
		panic("unreachable")
	}
}
