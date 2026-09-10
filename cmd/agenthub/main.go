package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wuyh266/agenthub/internal/agent"
	"github.com/wuyh266/agenthub/internal/llm"
	"github.com/wuyh266/agenthub/internal/tool"
	"github.com/wuyh266/agenthub/internal/tool/imageprocessing"
	"github.com/wuyh266/agenthub/internal/tool/mock"
)

func main() {
	apikey := os.Getenv("GLM_APIKEY")
	if apikey == "" {
		panic("GLM_APIKEY environment variable is not set")
	}
	llmClient := llm.NewGLMClient(apikey)
	tools := []tool.Tool{mock.Echo{}, imageprocessing.NewImageScaling("http://localhost:8080/resize")} 
	a := agent.NewAgent(llmClient, tools, 5)
	question:="使用Echo工具返回一个 hello world"
	if len(os.Args)>1{
		question = os.Args[1]
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	answer, err := a.Run(ctx, question)
	if err != nil {
		panic(err)
	}
	fmt.Println("最终答案:", answer)
}
