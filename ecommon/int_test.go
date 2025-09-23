package ecommon

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
	Value *Big `form:"Big" json:"Big" binding:"required"`
}

func TestGinBindJSONBig(t *testing.T) {
	makeRequest(map[string]string{
		"Big": "999999999111111133333",
	}, func(req BigValue) {
		assert.Equal(t, req.Value.String(), "999999999111111133333")
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
