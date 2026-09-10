package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wuyh266/agenthub/internal/agent"
	"github.com/wuyh266/agenthub/internal/llm"
	"github.com/wuyh266/agenthub/internal/tool"
	"github.com/wuyh266/agenthub/internal/tool/mock"
)

func main() {
	apikey := os.Getenv("GLM_APIKEY")
	if apikey == "" {
		panic("GLM_APIKEY environment variable is not set")
	}
	llmClient := llm.NewGLMClient(apikey)
	tools := []tool.Tool{mock.Echo{}}
	a := agent.NewAgent(llmClient, tools, 5)
	question := "请使用 echo 工具返回 'Hello, World!'"
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	answer, err := a.Run(ctx, question)
	if err != nil {
		panic(err)
	}
	fmt.Println("最终答案:", answer)
}
