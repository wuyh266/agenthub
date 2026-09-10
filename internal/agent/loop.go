package agent

import (
	"context"
	"errors"
	"strings"

	"github.com/wuyh266/agenthub/internal/llm"
	"github.com/wuyh266/agenthub/internal/tool"
)

type Agent struct {
	llm     llm.Client           // llm客户端，agent需要记住在跟谁说话
	tools   map[string]tool.Tool // 工具表，key是工具名，value是工具实例。agent需要记住有哪些工具可以调用
	maxLoop int                  // 最大循环次数。agent需要记住自己最多能循环多少次
}

// NewAgent 创建一个新的 Agent 实例。它接收一个 llm.Client、一个工具列表和最大循环次数作为参数，并返回一个指向 Agent 的指针。
func NewAgent(llm llm.Client, tools []tool.Tool, maxLoop int) *Agent {
	toolMap := make(map[string]tool.Tool)
	for _, t := range tools {
		toolMap[t.Name()] = t
	}
	return &Agent{
		llm:     llm,
		tools:   toolMap,
		maxLoop: maxLoop,
	}
}

func (a *Agent) buildPrompt() string {
	var prompt strings.Builder
	prompt.WriteString("你是一个智能体，能够调用工具来完成任务。你有以下工具可以使用：\n")
	for _, t := range a.tools {
		prompt.WriteString("- " + t.Name() + ": " + t.Description() + "\n")
	}
	prompt.WriteString("- 回复必须以 Thought: <你的思考> 开头\n")
	prompt.WriteString("- 需要用工具时，接着两行：Action: <必须是上面清单里的名字> 和 Action Input: <参数>，写完立即停止，等待 Observation\n")
	prompt.WriteString("- 能直接回答时，输出 Final Answer: <答案>\n")
	prompt.WriteString("- 每次回复只能调用一个工具\n")
	prompt.WriteString("现在，根据用户输入开始你的回复。\n")
	return prompt.String()
}

//  Run(ctx, 问题):
//         messages = [system提示词(工具清单+规则), user(问题)]
//         循环最多 maxLoop 次:
//             回复 = llm.Chat(ctx, messages)
//             把回复追加进 messages          // 模型说过的话要留档
//             decision = Parse(回复)
//             if decision.IsFinal: 返回 FinalAnswer
//             工具 = 查表(decision.ToolName)
//             观察结果 = 工具.Execute(ctx, decision.Input)
//             把 Observation 追加进 messages  //喂回去给模型看
//         循环耗尽: 返回错误

func (a *Agent) Run(ctx context.Context, question string) (string, error) {
	messages := []llm.Message{
		{Role: "system", Content: a.buildPrompt()},
		{Role: "user", Content: question},
	}
	for i := 0; i < a.maxLoop; i++ {
		reply, err := a.llm.Chat(ctx, messages)
		if err != nil {
			return "", err
		}
		messages = append(messages, llm.Message{Role: "assistant", Content: reply})
		decision := Parse(reply)
		if decision.IsFinal() {
			return decision.FinalAnswer, nil // 这个位置就直接返回了，结束循环
		}
		t, ok := a.tools[decision.ToolName] // 这里用t时因为tool会覆盖import tool，tools又不准确
		if !ok {
			messages = append(messages, llm.Message{Role: "user", Content: "Observation: 错误：工具 \"" + decision.ToolName + "\" 不存在"})
			continue
		}
		observation, err := t.Execute(ctx, decision.Input) // 这个observation是工具执行的结果，err是工具执行的错误
		if err != nil {
			messages = append(messages, llm.Message{Role: "user", Content: "Observation: 工具执行出错：" + err.Error()})
			continue
		}
		messages = append(messages, llm.Message{Role: "user", Content: "Observation: " + observation})
	}
	return "", errors.New("循环耗尽，未能得到最终答案")
}
