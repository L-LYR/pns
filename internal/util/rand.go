package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	_SourceBytes   = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_SourceIdxBits = 6
	_Mask          = 1<<_SourceIdxBits - 1
	_SourceIdxMax  = 63 / _SourceIdxBits
)

var _Rand = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
	sb := &strings.Builder{}
	sb.Grow(n)
	for i, cache, remain := n-1, _Rand.Int63(), _SourceIdxMax; i >= 0; remain-- {
		if remain == 0 {
			cache, remain = _Rand.Int63(), _SourceIdxMax
		}
		if idx := int(cache & _Mask); idx < len(_SourceBytes) {
			sb.WriteByte(_SourceBytes[idx])
			i--
		}
		cache >>= _SourceIdxBits
	}
	return sb.String()
}
