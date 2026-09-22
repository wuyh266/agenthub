package trace

import (
	"fmt"
	"io"
	"time"
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
	reply string,
	err error,
	elapsed time.Duration,
) {
	fmt.Fprintf(c.out, "\n=== Step %d ===\n", step)
	fmt.Fprintf(
		c.out,
		"模型调用耗时：%s\n",
		elapsed.Round(time.Millisecond),
	)

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
