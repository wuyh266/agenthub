package validation

import (
	"fmt"
	"github.com/santhosh-tekuri/jsonschema/v6"
	"strings"
)

const scaleSchemaJSON = `{
  "type": "object",
  "properties": {
    "path": {
      "type": "string"
    },
    "width": {
      "type": "integer",
      "minimum": 10,
      "maximum": 4096
    }
  },
  "required": ["path", "width"]
}`

func NewScaleValidator() (*jsonschema.Schema, error) {
	doc, err := jsonschema.UnmarshalJSON(strings.NewReader(scaleSchemaJSON)) //读入规则文档
	if err != nil {
		return nil, fmt.Errorf("读入规则文档失败: %w", err)
	}
	compiler := jsonschema.NewCompiler()                 //创建一个新的JSON schema编译器
	err = compiler.AddResource("scale.schema.json", doc) //注册规则文档
	if err != nil {
		return nil, fmt.Errorf("注册规则文档失败: %w", err)
	}
	compiledSchema, err := compiler.Compile("scale.schema.json") //编译规则文档
	if err != nil {
		return nil, fmt.Errorf("编译规则文档失败: %w", err)
	}
	return compiledSchema, nil

}
func ValidateScaleInput(schema *jsonschema.Schema, input string) error {
	getInput, err := jsonschema.UnmarshalJSON(strings.NewReader(input))
	if err != nil {
		return fmt.Errorf("解析输入JSON失败: %w", err)
	}
	err = schema.Validate(getInput)
	if err != nil {
		return fmt.Errorf("输入验证失败: %w", err)
	}
	return nil
}
