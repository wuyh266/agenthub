package tool

import "context"

// Tool 是一个工具接口。它定义了工具的基本行为，包括名称、描述和执行方法。
type Tool interface {
	// Name 返回工具名。LLM 在输出 Action 时用它点名，内核用它查表。必须唯一且不含空格。
	Name() string
	// Description 返回工具描述。内核把它拼进系统提示词，LLM 根据它决定调用哪个工具。它是工具的“说明书”，写给模型看的。
	Description() string
	// Execute 执行工具。调用者是内核，不是 LLM。LLM 只是“点菜”，内核拿着菜单指令调 Execute 干活，干完把结果端回去给 LLM 看。
	Execute(ctx context.Context, input string) (string, error)
}
