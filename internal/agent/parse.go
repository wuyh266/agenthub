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
	output = strings.ReplaceAll(output, "：", ":")
	decision := Decision{}
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		switch strings.Trim(strings.TrimSpace(key), "*") {
		case "Action":
			decision.ToolName = strings.TrimSpace(value)
		case "Action Input":
			inputLines := append(
				[]string{strings.TrimSpace(value)},
				lines[i+1:]...,
			)
			decision.Input = strings.TrimSpace(strings.Join(inputLines, "\n"))
			return decision
		case "Final Answer":
			answerLines := append(
				[]string{strings.TrimSpace(value)},
				lines[i+1:]...,
			)
			decision.FinalAnswer = strings.TrimSpace(strings.Join(answerLines, "\n"))
			return decision
		}
	}
	return decision
}
