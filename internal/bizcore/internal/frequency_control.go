package internal

import (
	"context"
	"time"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	log "github.com/L-LYR/pns/internal/service/push_log"
)

func _FreqCtrl(ctx context.Context, task model.PushTask) bool {
	if !task.GetMeta().UnderFreqCtrl() {
		return false
	}
	switch task.GetType() {
	case model.BroadcastPush:
		return _FreqCtrlByConfig(ctx, task, config.GetAppLevelFreqCtrlConfig())
	case model.DirectPush, model.RangePush:
		return _FreqCtrlByConfig(ctx, task, config.GetTargetLevelFreqCtrlConfig())
	default:
		return false
	}
}

func _FreqCtrlByConfig(
	ctx context.Context,
	task model.PushTask,
	cfg *config.FreqCtrlConfig,
) bool {
	// NOTICE: this counter will include both success and failure
	n, err := log.CountLogEntry(ctx, task.GetLogMeta(), cfg.Interval*time.Second)
	return err != nil || n > cfg.Limit
}

var (
	_FreqCtrlRule = `
rule "freqCtrl" "frequency control" salience 254
begin
	meta := task.GetMeta()
	if FrequencyControl(ctx, task) {
		meta.SetFiltered()
		stag.StopTag = true
	}
end
`
)
