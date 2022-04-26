package internal

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

// salience 0
// must be called at the end
func _EndObserve(ctx context.Context, task model.PushTask) {
	// NOTICE: Reserve for something special
}

const (
	_EndRule = `
rule "end" "observe at the end of validation" salience 0
begin
	EndObserve(ctx, task)
end	
`
)
