package llm

import (
	"context"
)

// Message 是一个消息类型。是对话中的一条消息，是发给模型的消息列表的基本单元
type Message struct {
	Role    string // 角色，分为 system、user、assistant 三种。模型生成的是 assistant，内核/工具生成的是 user，提示词是 system。
	Content string // 内容，文本内容。可以是任意文本，通常是自然语言。
}

// Client 是一个语言模型客户端接口。
type Client interface {
	Chat(ctx context.Context, messages []Message) (string, error)
}
