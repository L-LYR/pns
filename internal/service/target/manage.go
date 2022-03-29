package target

import (
	"context"
	"time"

	"github.com/L-LYR/pns/internal/local_storage"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/L-LYR/pns/internal/service/internal/do"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(ctx context.Context, target *model.Target) error {
	if err := dao.Target.Transaction(
		ctx,
		func(ctx context.Context, tx *gdb.TX) error {
			appName, _ := local_storage.GetAppNameByAppId(target.App.ID)

			if result, err := dao.TargetMongoDao.GetTarget(
				ctx,
				target.Device.ID,
				appName,
			); err != nil {
				return err
			} else if result != nil {
				EmitUpsertTargetEvent("create", "duplicate")
				util.GLog.Warningf(ctx, "%+v %+v", target.Device, target.App)
				return nil
			}

			now := time.Now()
			target.CreateTime = now
			target.LastActiveTime = now
			target.InfoUpdateTime = now

			if err := dao.TargetMongoDao.SetTarget(
				ctx,
				appName,
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
		util.GLog.Errorf(ctx, "%+v", err.Error())
		EmitUpsertTargetEvent("create", "failure")
		return err
	}
	return nil
}

func Update(ctx context.Context, target *model.Target) error {
	if err := dao.Target.Transaction(
		ctx,
		func(ctx context.Context, tx *gdb.TX) error {
			appName, _ := local_storage.GetAppNameByAppId(target.App.ID)

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

			if err := dao.TargetMongoDao.SetTarget(ctx, appName, target); err != nil {
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
		util.GLog.Errorf(ctx, "%+v", err.Error())
		EmitUpsertTargetEvent("update", "failure")
		return err
	}
	EmitUpsertTargetEvent("update", "success")
	return nil
}

func Query(ctx context.Context, deviceId string, appName string) (*model.Target, error) {
	return dao.TargetMongoDao.GetTarget(ctx, deviceId, appName)
}
