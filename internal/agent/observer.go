package agent

import (
	"time"
)

type Observer interface {
	OnLLMCall(step int, reply string, err error, elapsed time.Duration)
	OnToolCall(step int, name string, input string, output string, err error, elapsed time.Duration)
}
type NopObserver struct{}

func (NopObserver) OnLLMCall(step int, reply string, err error, elapsed time.Duration) {}

func (NopObserver) OnToolCall(step int, name string, input string, output string, err error, elapsed time.Duration) {
}
