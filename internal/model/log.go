package model

import (
	"fmt"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// To make the length of raw log shorter, we make each tag shorter
type PushLogMeta struct {
	TaskId   int
	AppId    int
	DeviceId string
}

func (b *PushLogMeta) EntryKey() string {
	return fmt.Sprintf("%d:%s", b.AppId, b.DeviceId)
}

func (b *PushLogMeta) StatusKey() string {
	return strconv.FormatInt(int64(b.TaskId), 10)
}

type LogBase struct {
	Meta  *PushLogMeta `json:"-"`
	T     int64        `json:"ts"`
	Where string       `json:"w"`
}

func (l *LogBase) Timestamp() int64 {
	return l.T
}

type LogEntry struct {
	*LogBase
	Hint string `json:"h"`
}

var (
	DummyEntry = &LogEntry{Hint: "Dummy"}
)

func (l *LogEntry) Decode(s string) error {
	return jsoniter.UnmarshalFromString(s, l)
}

func (l *LogEntry) Encode() (string, error) {
	return jsoniter.MarshalToString(l)
}

func (l *LogEntry) Readable() string {
	return fmt.Sprintf("[%s] %s at %s", time.UnixMilli(l.T).Format(time.RFC3339), l.Hint, l.Where)
}

// fyne notification does not provide onClick method
// so we cannot get the click event.