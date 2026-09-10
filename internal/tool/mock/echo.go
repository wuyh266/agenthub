package mock

import (
	"context"
)

type Echo struct{}

func (e Echo) Name() string {
	return "echo"
}

func (e Echo) Description() string {
	return "echo 工具。原样返回输入文本（附加 echo: 前缀）。输入为任意文本。适用于测试工具调用流程是否正常工作。"
}

func (e Echo) Execute(ctx context.Context, input string) (string, error) {
	return "echo" + input, nil
}
