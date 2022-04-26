package bizrule

import (
	"context"
	"errors"

	"github.com/L-LYR/pns/internal/bizcore"
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

func ChangeRuleStatus(
	ctx context.Context,
	name string,
	status model.BizRuleStatus,
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
			// update engine pool
			switch status {
			case model.BizRuleEnable:
				err = bizcore.AddNewRule(ctx, rule)
			case model.BizRuleDisable:
				err = bizcore.RemoveRule(ctx, name)
			}
			if err != nil {
				return err
			}
			return nil
		},
	)
}
