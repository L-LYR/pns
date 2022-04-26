package bizcore

import (
	"context"

	"github.com/L-LYR/pns/internal/bizcore/internal"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/util"
	"github.com/bilibili/gengine/engine"
)

var (
	_EnginePool *engine.GenginePool
)

type ExecuteMode int

const (
	Sorted ExecuteMode = 1 // by default
	// TODO: try these modes
	Concurrent ExecuteMode = 2
	Mixed      ExecuteMode = 3
	inverse    ExecuteMode = 4
)

func MustInitialize(ctx context.Context) {
	var err error
	// TODO: configurable
	_EnginePool, err = engine.NewGenginePool(
		config.GetEnginePoolMinLen(),
		config.GetEnginePoolMaxLen(),
		int(Sorted),
		internal.PredefinedRules(),
		internal.OuterApis(),
	)
	if err != nil {
		util.GLog.Panicf(ctx, "Fail to initialize engine pool, because %s", err.Error())
	}
}

func AddNewRule(ctx context.Context, rule *model.BizRule) error {
	return _EnginePool.UpdatePooledRulesIncremental(rule.String())
}

func RemoveRule(ctx context.Context, ruleName string) error {
	return _EnginePool.RemoveRules([]string{ruleName})
}

// NOTICE: I don't want to handle the return value...
func Execute(ctx context.Context, task model.PushTask) error {
	stag := &engine.Stag{}
	err, _ := _EnginePool.ExecuteWithStopTagDirect(
		map[string]interface{}{
			"stag": stag,
			"ctx":  ctx,
			"task": task,
		},
		false, stag)
	if err != nil {
		return err
	}
	return nil // well done
}
