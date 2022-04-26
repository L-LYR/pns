package internal

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

// salience 255
// must be called at the beginning
func _BeginObserve(ctx context.Context, task model.PushTask) {
	// NOTICE: Reserve for something special
}

const (
	_BeginRule = `
rule "begin" "observe at the begining of validation" salience 255
begin
	BeginObserve(ctx, task)
end
`
)
