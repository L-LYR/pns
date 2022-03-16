package target

import (
	"context"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/L-LYR/pns/internal/service/internal/do"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(ctx context.Context, target *model.Target) error {
	if err := dao.Target.Transaction(
		ctx,
		func(ctx context.Context, tx *gdb.TX) error {
			if result, err := dao.TargetMongoDao.GetTarget(
				ctx,
				target.Device.ID,
				target.App.ID,
			); err != nil {
				return err
			} else if result != nil {
				EmitUpsertTargetEvent("create", "duplicate")
				g.Log().Line().Warningf(ctx, "%+v", target)
				return nil
			}

			now := time.Now()
			target.CreateTime = now
			target.LastActiveTime = now
			target.InfoUpdateTime = now

			if err := dao.TargetMongoDao.SetTarget(
				ctx,
				target,
				options.Update().SetUpsert(true),
			); err != nil {
				return err
			}

			data := &do.Target{}
			if err := copier.Copy(data, target.Device); err != nil {
				return err
			}
			if err := copier.Copy(data, target.App); err != nil {
				return err
			}
			if _, err := dao.Target.Ctx(ctx).Data(data).Insert(); err != nil {
				return err
			}

			EmitUpsertTargetEvent("create", "success")
			return nil
		},
	); err != nil {
		g.Log().Line().Errorf(ctx, "%+v", err)
		EmitUpsertTargetEvent("create", "failure")
		return err
	}
	return nil
}

func Update(ctx context.Context, target *model.Target) error {
	if err := dao.Target.Transaction(
		ctx,
		func(ctx context.Context, tx *gdb.TX) error {
			updateData := do.Target{}
			if err := copier.Copy(&updateData, target); err != nil {
				return err
			}

			now := time.Now()
			target.LastActiveTime = now
			target.InfoUpdateTime = now
			if target.Tokens != nil && len(target.Tokens) > 0 {
				target.TokenUpdateTime = now
			}

			if err := dao.TargetMongoDao.SetTarget(ctx, target); err != nil {
				return err
			}

			if _, err := dao.Target.Ctx(ctx).
				Where("device_id", target.Device.ID).
				Where("app_id", target.App.ID).
				Update(updateData); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		g.Log().Line().Errorf(ctx, "%+v", err)
		EmitUpsertTargetEvent("update", "failure")
		return err
	}
	EmitUpsertTargetEvent("update", "success")
	return nil
}

func Query(ctx context.Context, deviceId string, appId int) (*model.Target, error) {
	return dao.TargetMongoDao.GetTarget(ctx, deviceId, appId)
}
