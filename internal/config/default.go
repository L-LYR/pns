package config

import (
	"context"

	"github.com/L-LYR/pns/internal/util"
	jsoniter "github.com/json-iterator/go"
)

func DefaultConfig(ctx context.Context) *Config {
	c := &Config{}
	if err := jsoniter.UnmarshalFromString(_RawDefaultConfig, c); err != nil {
		util.GLog.Panicf(ctx, "Default config is invalid, because %s", err.Error())
	}
	return c
}

var _RawDefaultConfig = `{
    "Servers": {
        "Inbound": {
            "Name": "inbound",
            "Address": ":10086"
        },
        "Bizapi": {
            "Name": "bizapi",
            "Address": ":10087"
        },
        "Admin": {
            "Name": "admin",
            "Address": ":10088"
        }
    },
    "Logger": {
        "Path": "./log",
        "File": "{Y-m-d}.log"
    },
    "Database": {
        "Mysql": {
            "Host": "localhost",
            "Port": "3306",
            "User": "root",
            "Pass": "pns_root",
            "Name": "pns"
        },
        "Mongo": {
            "User": "root",
            "Pass": "pns_root",
            "Host": "localhost",
            "Port": "27017",
            "Name": "pns_target"
        },
        "Redis": {
            "Pass": "pns_root",
            "Host": "localhost",
            "Port": "6379"
        }
    },
    "Broker": {
        "Broker": "localhost",
        "Port": "1883",
        "Timeout": 1
    },
    "EventQueue": {
        "TaskValidationEventQueue": {
            "Topic": "task_validation_event",
            "Dispatch": 5,
            "Length": 100000
        },
        "DirectPushTaskEventQueue": {
            "Topic": "direct_push_task_event",
            "Dispatch": 5,
            "Length": 100000
        },
        "RangePushTaskEventQueue": {
            "Topic": "range_push_task_event",
            "Dispatch": 5,
            "Length": 100000
        },
        "BroadcastPushTaskEventQueue": {
            "Topic": "broadcast_push_task_event",
            "Dispatch": 1,
            "Length": 10
        },
        "PushLogEventqueue": {
            "Topic": "push_log_event",
            "Dispatch": 3,
            "Length": 10000
        }
    },
    "EnginePool": {
        "MinLen": 10,
        "MaxLen": 20
    },
    "Misc": {
        "Qos": "atLeastOnce",
        "AuthKeyLength": 16,
        "AuthSecretLength": 32,
        "LogExpireTime": 86400,
        "TokenExpireTime": 604800,
        "MessageTemplateCacheSize": 100
    },
    "FrequencyControl": {
        "AppLevel": {
            "Interval": 3600,
            "Limit": 1
        },
        "TargetLevel": {
            "Interval": 3600,
            "Limit": 3
        }
    }
}`
