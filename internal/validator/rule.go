package validator

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/service/cache"
	"github.com/gogf/gf/v2/util/gvalid"
)

/*
	When parameters becomes more and more, they becomes to depend on each other.
	So here, we use rules to control the validation.
*/

type Rule struct {
	Name string
	Fn   gvalid.RuleFunc
}

var (
	_Rules = []*Rule{
		{
			Name: "app-exist",
			Fn: func(ctx context.Context, in gvalid.RuleFuncInput) error {
				if in.Value == nil {
					return errors.New("appid is empty")
				}
				ok := cache.Config.CheckAppExistByAppId(in.Value.Int())
				if ok {
					return nil
				}
				return errors.New("unknown app")
			},
		},
		{
			Name: "app-not-exist",
			Fn: func(ctx context.Context, in gvalid.RuleFuncInput) error {
				if in.Value == nil {
					return errors.New("appid is empty")
				}
				ok := cache.Config.CheckAppExistByAppId(in.Value.Int())
				if !ok {
					return nil
				}
				return errors.New("existed app")
			},
		},
	}
)

func MustRegisterRules(ctx context.Context) {
	for _, rule := range _Rules {
		gvalid.RegisterRule(rule.Name, rule.Fn)
	}
}
