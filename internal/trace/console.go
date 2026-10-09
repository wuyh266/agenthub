package trace

import (
	"fmt"
	"io"
	"time"

	"github.com/wuyh266/agenthub/internal/llm"
)

type ConsoleObserver struct {
	out io.Writer
}

func NewConsoleObserver(out io.Writer) *ConsoleObserver {
	return &ConsoleObserver{
		out: out,
	}
}
func (c *ConsoleObserver) OnLLMCall(
	step int,
	messageCount int,
	inputChars int,
	reply string,
	usage *llm.Usage,
	err error,
	elapsed time.Duration,
) {
	fmt.Fprintf(c.out, "\n=== Step %d ===\n", step)
	fmt.Fprintf(
		c.out,
		"模型输入：%d 条消息，%d 个字符\n",
		messageCount,
		inputChars,
	)
	fmt.Fprintf(
		c.out,
		"模型调用耗时：%s\n",
		elapsed.Round(time.Millisecond),
	)
	if usage == nil {
		fmt.Fprintln(c.out, "Token 用量：未取得")
	} else {
		fmt.Fprintf(
			c.out,
			"Token 用量：输入 %d，输出 %d，总计 %d\n",
			usage.PromptTokens,
			usage.CompletionTokens,
			usage.TotalTokens,
		)
	}

	if err != nil {
		fmt.Fprintf(c.out, "模型调用失败：%v\n", err)
		return
	}

	fmt.Fprintln(c.out, "模型回复：")
	fmt.Fprintln(c.out, reply)
}

func (c *ConsoleObserver) OnToolCall(
	step int,
	name string,
	input string,
	output string,
	err error,
	elapsed time.Duration,
) {
	fmt.Fprintf(
		c.out,
		"\n工具调用：%s，耗时：%s\n",
		name,
		elapsed.Round(time.Millisecond),
	)
	fmt.Fprintf(c.out, "工具输入：%s\n", input)

	if err != nil {
		fmt.Fprintf(c.out, "工具调用失败：%v\n", err)
		return
	}

	fmt.Fprintf(c.out, "工具结果：%s\n", output)
}

func (c *ConsoleObserver) OnRunEnd(
	totalUsage llm.Usage,
	llmCalls int,
	usageReportedCalls int,
) {
	fmt.Fprintln(c.out, "\n=== 本次任务统计 ===")
	fmt.Fprintf(
		c.out,
		"模型调用：%d 次，取得用量：%d 次\n",
		llmCalls,
		usageReportedCalls,
	)

	if usageReportedCalls == 0 {
		fmt.Fprintln(c.out, "累计 Token 用量：未取得")
		return
	}

	fmt.Fprintf(
		c.out,
		"累计 Token 用量：输入 %d，输出 %d，总计 %d\n",
		totalUsage.PromptTokens,
		totalUsage.CompletionTokens,
		totalUsage.TotalTokens,
	)

	if usageReportedCalls < llmCalls {
		fmt.Fprintln(c.out, "部分调用未返回用量，以上只累计已取得的用量。")
	}
}
