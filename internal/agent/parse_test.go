package agent

import "testing"

func TestParseActionWithFullWidthColon(t *testing.T) {
	output := "Action：echo"

	got := Parse(output)

	if got.ToolName != "echo" {
		t.Fatalf("ToolName = %q, want %q", got.ToolName, "echo")
	}
}

func TestParseActionWithMarkdownBold(t *testing.T) {
	output := "**Action**: echo"

	got := Parse(output)

	if got.ToolName != "echo" {
		t.Fatalf("ToolName = %q, want %q", got.ToolName, "echo")
	}
}
func TestParseMultilineActionInput(t *testing.T) {
	output := "Action: echo\nAction Input:\nhello world"

	got := Parse(output)

	if got.Input != "hello world" {
		t.Fatalf("Input = %q, want %q", got.Input, "hello world")
	}
}

func TestParseMultilineFinalAnswer(t *testing.T) {
	output := "Final Answer: 第一行\n第二行"

	got := Parse(output)

	want := "第一行\n第二行"
	if got.FinalAnswer != want {
		t.Fatalf("FinalAnswer = %q, want %q", got.FinalAnswer, want)
	}
}
func TestParseActionInputWithMultipleLines(t *testing.T) {
	output := "Action: echo\nAction Input:\n第一行\n第二行"

	got := Parse(output)

	want := "第一行\n第二行"
	if got.Input != want {
		t.Fatalf("Input = %q, want %q", got.Input, want)
	}
}
func TestParseStopsAfterActionInput(t *testing.T) {
	output := "Action: echo\nAction Input: hello\nAction: image_scaling"

	got := Parse(output)

	if got.ToolName != "echo" {
		t.Fatalf("ToolName = %q, want %q", got.ToolName, "echo")
	}
}
func TestParseCompleteToolRequest(t *testing.T) {
	output := "Thought: 需要缩放图片\nAction: image_scaling\nAction Input: ./picture/window.jpg 256"

	got := Parse(output)

	if got.ToolName != "image_scaling" {
		t.Fatalf("ToolName = %q, want %q", got.ToolName, "image_scaling")
	}
	if got.Input != "./picture/window.jpg 256" {
		t.Fatalf("Input = %q, want %q", got.Input, "./picture/window.jpg 256")
	}
}
func TestParseCompleteFinalAnswer(t *testing.T) {
	output := "Thought: 这个问题不需要调用工具\nFinal Answer: 这是最终答案"

	got := Parse(output)

	if !got.IsFinal() {
		t.Fatal("IsFinal() = false, want true")
	}
	if got.FinalAnswer != "这是最终答案" {
		t.Fatalf("FinalAnswer = %q, want %q", got.FinalAnswer, "这是最终答案")
	}
}
