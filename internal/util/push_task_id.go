package util

import "github.com/godruoyi/go-snowflake"

var (
	_DefaultPushTaskIdGenerator = _MustNewSnowflakeIdGenerator()
)

type _PushTaskIdGenerator interface {
	Get() uint64 // should be thread-safe
}

var (
	_ (_PushTaskIdGenerator) = (*_SnowflakeIdGenerator)(nil)
)

type _SnowflakeIdGenerator struct{}

func (s *_SnowflakeIdGenerator) Get() uint64 { return snowflake.ID() }

func _MustNewSnowflakeIdGenerator() *_SnowflakeIdGenerator {
	snowflake.SetMachineID(1)
	return &_SnowflakeIdGenerator{}
}

func GeneratePushTaskId() uint64 {
	return _DefaultPushTaskIdGenerator.Get()
}
