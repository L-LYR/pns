package util

import "github.com/godruoyi/go-snowflake"

/*
NOTICE:
	This package may be a little bit wired.
	Actually, I need an id generator which supports namespace
	to distinguish services, but I have not found one.
	Snowflake is a quite good id generator.
*/

var (
	_DefaultPushTaskIdGenerator = _MustNewSnowflakeIdGenerator(1)
	_DefaultTemplateIdGenerator = _MustNewSnowflakeIdGenerator(2)
)

func GeneratePushTaskId() int64 {
	return _DefaultPushTaskIdGenerator.Get()
}

func GenerateTemplateId() int64 {
	return _DefaultTemplateIdGenerator.Get()
}

type IdGenerator interface {
	Get() int64 // should be thread-safe
}

var (
	_ (IdGenerator) = (*_SnowflakeIdGenerator)(nil)
)

type _SnowflakeIdGenerator struct{}

func (s *_SnowflakeIdGenerator) Get() int64 { return int64(snowflake.ID()) }

func _MustNewSnowflakeIdGenerator(id uint16) *_SnowflakeIdGenerator {
	snowflake.SetMachineID(id)
	return &_SnowflakeIdGenerator{}
}
