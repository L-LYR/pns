package main

import (
	"github.com/L-LYR/pns/mobile/demo/ui"
	"github.com/L-LYR/pns/mobile/push_sdk"
)

func main() {
	push_sdk.MustInitialize(
		push_sdk.MustNewConfigFromString(
			_RawSettings,
			push_sdk.GenerateUUIDAsDeviceId(),
		),
		ui.LogHandler,
	)
	ui.Run()
}
