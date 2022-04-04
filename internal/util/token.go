package util

import (
	"crypto/sha1"
	"fmt"
	"hash"
	"time"
)

type TokenBuilder interface {
	Build(*TokenSource) (string, error)
}

var (
	_ TokenBuilder = (*_NaiveSHA1Token)(nil)
)

func NewTokenBuilder() TokenBuilder {
	return &_NaiveSHA1Token{b: sha1.New()}
}

type TokenSource struct {
	AppId    int
	DeviceId string
}

// We suppose that DeviceId is different between devices.
// We use uuid as a temporary solution, which should not
// be considered in this project.
type _NaiveSHA1Token struct {
	b hash.Hash
}

func (t *_NaiveSHA1Token) Build(ts *TokenSource) (string, error) {
	t.b.Reset()
	s := fmt.Sprintf("%d-%s-%d", ts.AppId, ts.DeviceId, time.Now().UnixNano())
	if _, err := t.b.Write([]byte(s)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", t.b.Sum(nil)), nil
}
