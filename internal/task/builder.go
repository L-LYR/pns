package task

import (
	"context"
	"errors"
	"time"

	v1 "github.com/L-LYR/pns/internal/bizapi/api/v1"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/cache"
	"github.com/L-LYR/pns/internal/service/target"
	"github.com/L-LYR/pns/internal/util"
)

/*
First-Of-All:
	specify task type

Common Procedures
	1. GenerateID
	2. SetMeta(Qos Pusher Retry CreateTime)

Normal Message:
	1. just set it

Template Message:
	1. fill in template

Direct Push:
	1. query target

Broadcast Push:
	1. set app
*/

type TaskBuilder interface {
	SetTaskMeta(int) TaskBuilder
	SetMessage(*v1.BasicMessage) TaskBuilder
	SetTemplateMessage(*v1.TemplateMessage) TaskBuilder
	SetDirectPushBase(*v1.DirectPushBase) TaskBuilder
	SetBroadcastPushBase(*v1.BroadcastPushBase) TaskBuilder

	Build() (model.PushTask, error)
}

func NewTaskBuilder(ctx context.Context, t model.PushTaskType) TaskBuilder {
	switch t {
	case model.DirectPush:
		return &DirectPushTaskBuilder{
			ctx:  ctx,
			task: &model.DirectPushTask{},
			err:  nil,
		}
	case model.BroadcastPush:
		return &BroadcastPushTaskBuilder{
			ctx:  ctx,
			task: &model.BroadcastTask{},
			err:  nil,
		}
	default:
		panic("unreachable")
	}
}

var (
	_ (TaskBuilder) = (*DirectPushTaskBuilder)(nil)
	_ (TaskBuilder) = (*BroadcastPushTaskBuilder)(nil)
)

type DirectPushTaskBuilder struct {
	ctx  context.Context
	task *model.DirectPushTask
	err  error
}

func (b *DirectPushTaskBuilder) SetTaskMeta(retry int) TaskBuilder {
	b.task.PushTaskMeta = _NewTaskMeta(retry)
	return b
}

func (b *DirectPushTaskBuilder) SetMessage(m *v1.BasicMessage) TaskBuilder {
	b.task.Message = _NewMessage(m)
	return b
}

func (b *DirectPushTaskBuilder) SetTemplateMessage(m *v1.TemplateMessage) TaskBuilder {
	b.task.Message, b.err = _NewMessageFromTemplate(b.ctx, m)
	return b
}

func (b *DirectPushTaskBuilder) SetDirectPushBase(base *v1.DirectPushBase) TaskBuilder {
	target, err := target.Query(b.ctx, base.AppId, base.DeviceId)
	if err != nil {
		b.err = errors.New("fail to query target")
		return b
	}
	if target == nil {
		b.err = errors.New("target not found")
		return b
	}
	b.task.Target = target
	return b
}

func (b *DirectPushTaskBuilder) SetBroadcastPushBase(*v1.BroadcastPushBase) TaskBuilder {
	return b
}

func (b *DirectPushTaskBuilder) Build() (model.PushTask, error) {
	return b.task, b.err
}

type BroadcastPushTaskBuilder struct {
	ctx  context.Context
	task *model.BroadcastTask
	err  error
}

func (b *BroadcastPushTaskBuilder) SetTaskMeta(retry int) TaskBuilder {
	b.task.PushTaskMeta = _NewTaskMeta(retry)
	return b
}

func (b *BroadcastPushTaskBuilder) SetMessage(m *v1.BasicMessage) TaskBuilder {
	b.task.Message = _NewMessage(m)
	return b
}

func (b *BroadcastPushTaskBuilder) SetTemplateMessage(m *v1.TemplateMessage) TaskBuilder {
	b.task.Message, b.err = _NewMessageFromTemplate(b.ctx, m)
	return b
}

func (b *BroadcastPushTaskBuilder) SetDirectPushBase(*v1.DirectPushBase) TaskBuilder {
	return b
}

func (b *BroadcastPushTaskBuilder) SetBroadcastPushBase(base *v1.BroadcastPushBase) TaskBuilder {
	b.task.AppId = base.AppId
	return b
}

func (b *BroadcastPushTaskBuilder) Build() (model.PushTask, error) {
	return b.task, b.err
}

func _NewTaskMeta(retry int) *model.PushTaskMeta {
	return &model.PushTaskMeta{
		RetryCounter: &model.RetryCounter{
			Counter: model.RetryTimes(retry),
		},
		ID:             util.GeneratePushTaskId(),
		Pusher:         model.MQTTPusher,
		Qos:            config.CommonTaskQos(),
		Status:         model.Pending,
		CreationTime:   time.Now(),
		IgnoreFreqCtrl: false,
	}
}

func _NewMessage(msg *v1.BasicMessage) *model.Message {
	return &model.Message{
		Title:   msg.Title,
		Content: msg.Content,
	}
}

func _NewMessageFromTemplate(ctx context.Context, tplMsg *v1.TemplateMessage) (*model.Message, error) {
	tpl, err := cache.MsgTpl.GetTplByID(ctx, tplMsg.Id)
	if err != nil {
		return nil, err
	}
	replaceStrings := make(map[string]map[string]string)
	for field, params := range tplMsg.ParamLists {
		replaceStrings[field] = params.PR
	}
	return tpl.FillInParams(replaceStrings)
}
