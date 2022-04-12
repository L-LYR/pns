package util

import "github.com/godruoyi/go-snowflake"

var (
	_DefaultPushTaskIdGenerator = _MustNewSnowflakeIdGenerator()
)

type _PushTaskIdGenerator interface {
	Get() int // should be thread-safe
}

var (
	_ (_PushTaskIdGenerator) = (*_SnowflakeIdGenerator)(nil)
)

type _SnowflakeIdGenerator struct{}

func (s *_SnowflakeIdGenerator) Get() int { return int(snowflake.ID()) }

func _MustNewSnowflakeIdGenerator() *_SnowflakeIdGenerator {
	snowflake.SetMachineID(1)
	return &_SnowflakeIdGenerator{}
}

func GeneratePushTaskId() int {
	return _DefaultPushTaskIdGenerator.Get()
}
