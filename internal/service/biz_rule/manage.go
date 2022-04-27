package bizrule

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/L-LYR/pns/internal/service/internal/do"
	"github.com/gogf/gf/v2/database/gdb"
)

func InsertBizRule(
	ctx context.Context,
	name string,
	description string,
	salience int,
	content string,
) error {
	rule := &model.BizRule{
		Name:        name,
		Description: description,
		Salience:    salience,
		Content:     content,
		Status:      model.BizRuleDisable,
	}
	if err := dao.InsertBizRule(ctx, rule); err != nil {
		return err
	}
	return nil
}

func LoadAllRules(ctx context.Context) ([]*model.BizRule, error) {
	return dao.LoadAllRules(ctx)
}

func ChangeRuleStatus(
	ctx context.Context,
	name string,
	status model.BizRuleStatus,
	fn func(context.Context, *model.BizRule) error,
) error {
	return dao.BizRule.Transaction(ctx,
		func(ctx context.Context, tx *gdb.TX) error {
			// query rule
			rule, err := dao.QueryBizRule(ctx, name)
			if err != nil {
				return err
			}
			if rule == nil {
				return errors.New("not found")
			}
			if rule.Status == status {
				return nil
			}
			// update database
			if res, err := dao.BizRule.Ctx(ctx).Data(do.BizRule{
				Status: status,
			}).Where("name", name).Update(); err != nil {
				return err
			} else if n, err := res.RowsAffected(); err != nil || n == 0 {
				return errors.New("unchanged")
			}
			// do something else in txn
			return fn(ctx, rule)
		},
	)
}
