package cache

import (
	"context"

	"github.com/L-LYR/pns/internal/config"
	"github.com/L-LYR/pns/internal/model"
	"github.com/L-LYR/pns/internal/service/internal/dao"
	"github.com/L-LYR/pns/internal/util"
	"github.com/gogf/gf/v2/os/gcache"
)

type _MsgTplCache struct {
	cache *gcache.Cache
}

func (c *_MsgTplCache) GetTplByID(ctx context.Context, id int64) (*model.MsgTpl, error) {
	if v, err := c.cache.Get(ctx, id); err != nil {
		return nil, err
	} else if v != nil { // cache miss
		msgTpl := &model.MsgTpl{}
		if err := v.Struct(msgTpl); err != nil {
			return nil, err
		}
		return msgTpl, nil
	} else {
		msgTpl, err := dao.QueryMessageTemplate(ctx, id)
		if err != nil {
			return nil, err
		}
		err = c.cache.Set(ctx, id, msgTpl, 0)
		if err != nil {
			util.GLog.Warningf(ctx, "Fail to set cache for template %d, because %s", id, err.Error())
		}
		return msgTpl, nil
	}
}

func (c *_MsgTplCache) MustInitialize(ctx context.Context) {
	c.cache = gcache.New(config.GetMsgTplCacheSize())
}

var (
	MsgTpl = &_MsgTplCache{}
)
