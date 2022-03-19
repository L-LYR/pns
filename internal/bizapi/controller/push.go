package controller

import (
	"context"

	"github.com/L-LYR/pns/internal/bizapi/api/v1"
)

var Push = _PushAPI{}

type _PushAPI struct{}

func (api *_PushAPI) Push(ctx context.Context, request *v1.PushReq) (*v1.PushRes, error) {
	response := &v1.PushRes{
		PushTaskId: "Hello World",
	}
	return response, nil
}
