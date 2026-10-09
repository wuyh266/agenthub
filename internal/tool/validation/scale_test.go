package validation

import (
	"errors"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

func TestScaleValidator(t *testing.T) {
	validator, err := NewScaleValidator()
	if err != nil {
		t.Fatalf("创建ScaleValidator失败: %v", err)
	}
	if validator == nil {
		t.Fatalf("ScaleValidator为nil")
	}
}

func TestValidateScaleInput(t *testing.T) {
	validator, err := NewScaleValidator()
	if err != nil {
		t.Fatalf("创建ScaleValidator失败: %v", err)
	}

	validInput := `{"path": "test/path", "width": 128}`
	err = ValidateScaleInput(validator, validInput)
	if err != nil {
		t.Errorf("验证有效输入失败: %v", err)
	}

	invalidInput := `{"path": "test/path", "width": "128"}`
	err = ValidateScaleInput(validator, invalidInput)
	if err == nil {
		t.Errorf("验证无效输入未返回错误")
	}
}

func TestValidateScaleInputRules(t *testing.T) {
	validator, err := NewScaleValidator()
	if err != nil {
		t.Fatalf("创建ScaleValidator失败: %v", err)
	}
	if validator == nil {
		t.Fatal("ScaleValidator为nil")
	}

	// wantError 为空表示通过；json 表示解析错误；schema 表示规则错误。
	tests := []struct {
		name      string
		input     string
		wantError string
	}{
		{"width_below_minimum", `{"path":"p.png","width":9}`, "schema"},
		{"width_at_minimum", `{"path":"p.png","width":10}`, ""},
		{"width_at_maximum", `{"path":"p.png","width":4096}`, ""},
		{"width_above_maximum", `{"path":"p.png","width":4097}`, "schema"},
		{"width_integral_decimal", `{"path":"p.png","width":128.0}`, ""},
		{"width_fractional", `{"path":"p.png","width":128.5}`, "schema"},
		{"width_string", `{"path":"p.png","width":"128"}`, "schema"},
		{"width_null", `{"path":"p.png","width":null}`, "schema"},
		{"missing_both_fields", `{}`, "schema"},
		{"missing_path", `{"width":128}`, "schema"},
		{"missing_width", `{"path":"p.png"}`, "schema"},
		{"root_null", `null`, "schema"},
		{"root_array", `[]`, "schema"},
		{"path_number", `{"path":128,"width":128}`, "schema"},
		{"path_null", `{"path":null,"width":128}`, "schema"},
		{"empty_path_allowed", `{"path":"","width":128}`, ""},
		{"extra_field_allowed", `{"path":"p.png","width":128,"format":"png"}`, ""},
		{"malformed_json", `{"path":"p.png","width":"128}`, "json"},
		{"multiple_json_values", `{"path":"p.png","width":128} {"path":"p.png","width":256}`, "json"},
		{"empty_input", ``, "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateScaleInput(validator, tt.input)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("期望校验通过，实际错误: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("期望 %s 错误，实际校验通过", tt.wantError)
			}

			// 检查包装后的底层错误类型，避免格式错误冒充规则校验结果。
			var validationErr *jsonschema.ValidationError
			isSchemaError := errors.As(err, &validationErr)
			if isSchemaError != (tt.wantError == "schema") {
				t.Fatalf("错误阶段不符，期望 %s 错误，实际错误: %v", tt.wantError, err)
			}
		})
	}
}
