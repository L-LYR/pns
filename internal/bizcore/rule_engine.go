package bizcore

import (
	"context"

	"github.com/L-LYR/pns/internal/bizcore/internal"
	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	bizrule "github.com/L-LYR/pns/internal/service/biz_rule"
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

	// Load Rules From Database

	rules, err := bizrule.LoadAllRules(ctx)
	if err != nil {
		util.GLog.Panicf(ctx, "Fail to load rules from database, because %s", err.Error())
	}
	for _, rule := range rules {
		if rule.Status == model.BizRuleDisable {
			continue
		}
		if err := AddNewRule(ctx, rule); err != nil {
			util.GLog.Panicf(ctx, "Fail to add rule %s, because %s", rule.Name, err.Error())
		}
	}
	util.GLog.Infof(ctx, "Load %d rules successfully", _EnginePool.GetRulesNumber())
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
