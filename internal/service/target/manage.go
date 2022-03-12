package target

import (
	"context"
	"time"

	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/monitor"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/L-LYR/pns/internal/service/internal/do"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/jinzhu/copier"
)

func Create(ctx context.Context, target *model.Target) error {
	g.Log().Line().Infof(ctx, "%+v %+v", target.Device, target.App)
	if err := dao.Target.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		have, err := dao.Target.Ctx(ctx).
			Where("device_id", target.Device.ID).
			Where("app_id", target.App.ID).
			Count()
		if err != nil {
			return err
		}
		if have > 0 {
			monitor.UpsertTargetMetric.WithLabelValues("create", "duplicate").Inc()
			g.Log().Line().Warningf(ctx, "%+v", target)
			return nil
		}

		data := &do.Target{}
		if err := copier.Copy(data, target.Device); err != nil {
			return err
		}
		if err := copier.Copy(data, target.App); err != nil {
			return err
		}
		data.CreateTime = time.Now()
		if _, err := dao.Target.Ctx(ctx).Data(data).Insert(); err != nil {
			return err
		}
		monitor.UpsertTargetMetric.WithLabelValues("create", "success").Inc()
		return nil
	}); err != nil {
		g.Log().Line().Errorf(ctx, "%+v", err)
		monitor.UpsertTargetMetric.WithLabelValues("create", "failure").Inc()
		return err
	}
	return nil
}

func Update(ctx context.Context, target *model.Target) error {
	g.Log().Line().Infof(ctx, "%+v", target)
	if err := dao.Target.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		updateData := do.Target{}
		if err := copier.Copy(&updateData, target); err != nil {
			return err
		}

		if _, err := dao.Target.Ctx(ctx).
			Where("device_id", target.Device.ID).
			Where("app_id", target.App.ID).
			Update(updateData); err != nil {
			return err
		}
		return nil
	}); err != nil {
		g.Log().Line().Errorf(ctx, "%+v", err)
		monitor.UpsertTargetMetric.WithLabelValues("update", "failure").Inc()
		return err
	}
	monitor.UpsertTargetMetric.WithLabelValues("update", "success").Inc()
	return nil
}
