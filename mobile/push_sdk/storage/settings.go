package storage

type Settings struct {
	Sdk SDKSettings
	App APPSettings
}

type SDKSettings struct {
	Version string
	MQTT    MQTTSettings
	Inbound InboundSettings
}

type MQTTSettings struct {
	Broker         string
	Port           string
	RetryInterval  int64
	ConnectTimeout int64
}

type InboundSettings struct {
	Base string
}

type APPSettings struct {
	ID      int
	Key     string
	Secret  string
	Version string
}
