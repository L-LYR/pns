package config

import jsoniter "github.com/json-iterator/go"

func DefaultConfig() *Config {
	c := &Config{}
	if err := jsoniter.UnmarshalFromString(_RawDefaultConfig, c); err != nil {
		panic("default config is invalid")
	}
	return c
}

var _RawDefaultConfig = `{
    "Servers": {
        "Inbound": {
            "Address": ":10086"
        },
        "Bizapi": {
            "Address": ":10087"
        },
        "Admin": {
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
        "PushTaskEventConsumer": {
            "Topic": "push_task_event",
            "Dispatch": 5
        },
        "PushLogEventConsumer": {
            "Topic": "push_log_event",
            "Dispatch": 5
        }
    },
    "Misc": {
        "Qos": "atLeastOnce"
    }
}`
