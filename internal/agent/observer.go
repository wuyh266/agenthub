package agent

import (
	"time"

	"github.com/wuyh266/agenthub/internal/llm"
)

type Observer interface {
	OnLLMCall(
		step int,
		messageCount int,
		inputChars int,
		reply string,
		usage *llm.Usage,
		err error,
		elapsed time.Duration,
	) //改的时候注意看，这里是接口
	OnToolCall(
		step int,
		name string,
		input string,
		output string,
		err error,
		elapsed time.Duration,
	) //这里也是接口！
	OnRunEnd(
		totalUsage llm.Usage,
		llmCalls int,
		usageReportedCalls int,
	)
}

type NopObserver struct{}

func (NopObserver) OnLLMCall(step int, messageCount int, inputChars int, reply string, usage *llm.Usage, err error, elapsed time.Duration) {
}

func (NopObserver) OnToolCall(step int, name string, input string, output string, err error, elapsed time.Duration) {
}

func (NopObserver) OnRunEnd(totalUsage llm.Usage, llmCalls int, usageReportedCalls int) {}
