package bizcore

import (
	"context"

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
	_EnginePool, err = engine.NewGenginePool(10, 20, int(Sorted), _DefaultRule, _ApiOuter)
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
func Execute(data map[string]interface{}) error {
	stag := &engine.Stag{}
	data["stag"] = stag
	err, _ := _EnginePool.ExecuteWithStopTagDirect(data, false, stag)
	if err != nil {
		return err
	}
	return nil // well done
}
