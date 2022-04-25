package util

import "github.com/godruoyi/go-snowflake"

var (
	_DefaultPushTaskIdGenerator = _MustNewSnowflakeIdGenerator()
)

type _PushTaskIdGenerator interface {
	Get() int64 // should be thread-safe
}

var (
	_ (_PushTaskIdGenerator) = (*_SnowflakeIdGenerator)(nil)
)

type _SnowflakeIdGenerator struct{}

func (s *_SnowflakeIdGenerator) Get() int64 { return int64(snowflake.ID()) }

func _MustNewSnowflakeIdGenerator() *_SnowflakeIdGenerator {
	snowflake.SetMachineID(1)
	return &_SnowflakeIdGenerator{}
}

func GeneratePushTaskId() int64 {
	return _DefaultPushTaskIdGenerator.Get()
}
