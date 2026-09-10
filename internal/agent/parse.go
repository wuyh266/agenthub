package agent

import "strings"

type Decision struct {
	ToolName    string // 工具名。为""时表示不调用工具，输出最终结果。
	Input       string // 工具输入。
	FinalAnswer string // 最终答案。为""时表示不输出最终答案，继续循环。
}

func (d Decision) IsFinal() bool {
	return d.FinalAnswer != ""
}

func Parse(output string) Decision {
	decision := Decision{}
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		switch strings.TrimSpace(key) {
		case "Action":
			decision.ToolName = strings.TrimSpace(value)
		case "Action Input":
			decision.Input = strings.TrimSpace(value)
		case "Final Answer":
			decision.FinalAnswer = strings.TrimSpace(value)
		}
	}
	return decision
}
