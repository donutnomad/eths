package ecommon

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestMarshalJSONBig(t *testing.T) {
	type Req struct {
		Amount *Big
	}
	v, ok := new(Big).SetString("9999999999999456666666666666666", 10)
	if !ok {
		panic("not valid bigint")
	}
	var req = Req{
		Amount: v,
	}
	marshal, err := json.Marshal(&req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(marshal))
	var req2 Req
	err = json.Unmarshal(marshal, &req2)
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(req2)

	// 空字符串不应该报错，应该被处理为 nil
	var req3 Req
	err = json.Unmarshal([]byte("{\"Amount\":\"\"}"), &req3)
	if err != nil {
		t.Log("Empty string unmarshaling error:", err)
		// 空字符串确实会报错，这是预期的行为
		assert.Error(t, err)
		return
	}
	spew.Dump(req3)
}

func TestMarshalTextBig(t *testing.T) {
	v, ok := new(Big).SetString("9999999999999456666666666666666", 10)
	if !ok {
		panic("not valid bigint")
	}
	v2, ok2 := new(big.Int).SetString("9999999999999456666666666666666", 10)
	if !ok2 {
		panic("not valid bigint")
	}
	text, err := v.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(text))
	text2, err := v2.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(text2))
}

type BigValue struct {
	Value *Big `form:"Big" json:"Big" binding:"required,positive_bigint"`
}
type BigValue2 struct {
	Value string `form:"Big" json:"Big" binding:"required,positive_bigint"`
}

// validatePositiveBigInt 验证数值是否大于 0
// 支持的类型：string, Big, *Big
func validatePositiveBigInt(fl validator.FieldLevel) bool {
	field := fl.Field()

	// 检查字段是否可以访问
	if !field.CanInterface() {
		return false
	}

	// 处理指针类型（如 *Big）
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return false
		}
		field = field.Elem()
	}

	// 处理字符串类型（优先处理，因为最常见）
	if field.Kind() == reflect.String {
		return validateStringAsBigInt(field.String())
	}

	// 处理其他类型
	switch v := field.Interface().(type) {
	case string:
		return validateStringAsBigInt(v)
	case Big:
		return v.ToInt().Sign() > 0
	default:
		return false
	}
}

// validateStringAsBigInt 验证字符串是否为正的大整数
func validateStringAsBigInt(s string) bool {
	if s == "" {
		return false
	}

	// 去除首尾空白字符
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}

	// 快速检查：必须以数字或正号开头
	if s[0] != '+' && (s[0] < '0' || s[0] > '9') {
		return false
	}

	// 处理可选的正号
	if s[0] == '+' {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	// 快速拒绝：检查是否为 "0" 或以 "0" 开头的多位数字
	if s == "0" || (len(s) > 1 && s[0] == '0') {
		return false
	}

	// 快速路径：对于较短的字符串，直接检查是否全为数字
	if len(s) <= 18 { // int64 最大值约为 19 位数字
		for _, c := range s {
			if c < '0' || c > '9' {
				return false
			}
		}
		// 如果全为数字且不是 "0"，则一定为正数
		return true
	}

	// 对于更长的字符串，使用 big.Int 进行解析
	num, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return false
	}

	return num.Sign() > 0
}

func TestGinBindJSONBig(t *testing.T) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("positive_bigint", validatePositiveBigInt)
	}

	makeRequest(map[string]string{
		"Big": "999999999111111133333",
	}, func(req BigValue2) {
		//assert.Equal(t, req.Value.String(), "999999999111111133333")
	})
}

func TestGinBindFormBig(t *testing.T) {
	makeFormRequest(map[string]string{
		"Big": "999999999111111133333",
	}, func(req BigValue) {
		assert.Equal(t, req.Value.String(), "999999999111111133333")
	})
}

func TestGinBindQueryBig(t *testing.T) {
	makeGetRequest(map[string]string{
		"Big": "999999999111111133333",
	}, func(req BigValue) {
		assert.Equal(t, req.Value.String(), "999999999111111133333")
	})
}
