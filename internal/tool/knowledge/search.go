package knowledge

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Search struct {
	path string
}

func NewSearch(path string) *Search {
	return &Search{
		path: path,
	}
}

func (s *Search) Name() string {
	return "knowledge_search"
}

func (s *Search) Description() string {
	return "检索演示商店的退货、保修和运费政策。" +
		"输入为用空格分隔的关键词，例如：耳机 拆封 退货。" +
		"回答商店政策问题时应先检索，依据返回资料回答并标注来源。" +
		"未找到相关资料时，应说明信息不足。"
}

func (s *Search) Execute(ctx context.Context, input string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	keywords := strings.Fields(strings.ToLower(input))
	if len(keywords) == 0 {
		return "", fmt.Errorf("检索关键词不能为空")
	}

	content, err := os.ReadFile(s.path)
	if err != nil {
		return "", fmt.Errorf("读取知识文档失败: %w", err)
	}

	sections := strings.Split(string(content), "\n## ")
	if len(sections) < 2 {
		return "", fmt.Errorf("知识文档缺少二级标题")
	}

	var matches []string
	for _, section := range sections[1:] {
		searchText := strings.ToLower(section)

		for _, keyword := range keywords {
			if strings.Contains(searchText, keyword) {
				matches = append(
					matches,
					"来源："+filepath.Base(s.path)+
						"\n## "+strings.TrimSpace(section),
				)
				break
			}
		}
	}

	if len(matches) == 0 {
		return "未找到相关资料。当前知识库无法提供这个问题的答案。", nil
	}

	return strings.Join(matches, "\n\n---\n\n"), nil
}
