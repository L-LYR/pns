package internal

/*
	NOTICE: this rule is just for testing
*/

import (
	"math/rand"
	"time"
)

func _RandomDelay() {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
}

var (
	_RandomDelayRule = `
rule "radomDelay" "random delay" salience 1
begin
	RandomDelay()
end
`
)
