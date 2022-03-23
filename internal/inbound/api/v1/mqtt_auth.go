package v1

import "github.com/gogf/gf/v2/frame/g"

type MQTTAuthReq struct {
	g.Meta   `path:"/client_auth" method:"post"`
	Username string `json:"username" v:"required#username is required"`
	Password string `json:"password" v:"required#password is required"`
	ClientId string `json:"clientid" v:"required#clientid is required"`
}

type MQTTAuthRes struct{ CommonAuthRes }

type ACLCheckReq struct {
	g.Meta   `path:"/acl_check" method:"post"`
	Username string `json:"username" v:"required#username is required"`
	ClientId string `json:"clientid" v:"required#clientid is required"`
	Topic    string `json:"topic" v:"required#topic is required"`
	Acc      string `json:"acc" v:"required#acc is required"`
}

type ACLCheckRes struct{ CommonAuthRes }

type CommonAuthRes struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}
