package bizcore

import (
	"context"

	"github.com/L-LYR/pns/internal/model"
)

/*
	Functions exposed to rules should be wrapperred or defined here
	and be registered in _ApiOuter.
	These functions are not exposed to any other packages, but only used in rules.
*/

var (
	_ApiOuter = map[string]interface{}{
		"BeginObserve": BeginObserve,
		"EndObserve":   EndObserve,
	}
)

// salience 255
// must be called at the beginning
func BeginObserve(ctx context.Context, task model.PushTask) {}

// salience 0
// must be called at the end
func EndObserve(ctx context.Context, task model.PushTask) {}

const (
	_DefaultRule = `
rule "begin" "observe at the begining of validation" salience 255
begin
	BeginObserve(ctx, task)
end	

rule "end" "observe at the end of validation" salience 0
begin
	EndObserve(ctx, task)
end	
`
)
