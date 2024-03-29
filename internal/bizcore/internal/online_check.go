package internal

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/target_status"
)

func _CheckTargetOnline(ctx context.Context, task model.PushTask) bool {
	if !task.GetMeta().NeedOnlineCheck() {
		return true
	}
	switch task.GetType() {
	case model.BroadcastPush, model.RangePush:
		return true
	case model.DirectPush:
		return target_status.CheckTargetOnline(
			ctx,
			task.GetAppId(),
			model.AsDirectPushTask(task).Device.ID,
		)
	default:
		return false
	}
}

var (
	_OnlineCheckRule = `
rule "onlineCheck" "online check" salience 253
begin
	meta := task.GetMeta()
	if !CheckTargetOnline(ctx, task) {
		meta.SetFiltered()
		stag.StopTag = true
	}
end
`
)
