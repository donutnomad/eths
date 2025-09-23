package ecommon

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type UnprefixedRequest struct {
	HashValue UnprefixedHash `json:"hash_value"`
}

type UnprefixedCompareRequest struct {
	Hash1 UnprefixedHash `json:"hash1" form:"hash1"`
	Hash2 UnprefixedHash `json:"hash2" form:"hash2"`
}

func TestUnprefixedHashJson(t *testing.T) {
	makeRequest(map[string]string{
		"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}, func(req UnprefixedCompareRequest) {
		cmp := req.Hash1.Cmp(req.Hash2)
		assert.Equal(t, 0, cmp)
		assert.True(t, cmp == 0)
		assert.Equal(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash1.Hex())
		assert.Equal(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash2.Hex())
	})
}

func TestUnprefixedHashForm(t *testing.T) {
	makeFormRequest(map[string]string{
		"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}, func(req UnprefixedCompareRequest) {
		cmp := req.Hash1.Cmp(req.Hash2)
		assert.Equal(t, 0, cmp)
		assert.True(t, cmp == 0)
		assert.Equal(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash1.Hex())
		assert.Equal(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash2.Hex())
	})
}

func TestUnprefixedHashFormFailed(t *testing.T) {
	assert.Panics(t, func() {
		makeFormRequest(map[string]string{
			"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefed",
			"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		}, func(req UnprefixedCompareRequest) {

		})
	})
}

// 正常
func TestUnprefixedHashQuery(t *testing.T) {
	makeGetRequest(map[string]string{
		"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
	}, func(req UnprefixedCompareRequest) {
		cmp := req.Hash1.Cmp(req.Hash2)
		assert.Equal(t, 0, cmp)
		assert.True(t, cmp == 0)
		assert.Equal(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash1.Hex())
		assert.Equal(t, "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash2.Hex())
	})
}

// 超过32个字节，报错
func TestUnprefixedHashQuery2(t *testing.T) {
	assert.Panics(t, func() {
		makeGetRequest(map[string]string{
			"hash1": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdefed",
			"hash2": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		}, func(req UnprefixedCompareRequest) {
			spew.Dump(req)
			//cmp := req.Hash1.Cmp(req.Hash2)
			//assert.Equal(t, 0, cmp)
			//assert.True(t, cmp == 0)
			//assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash1.Hex())
			//assert.Equal(t, "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", req.Hash2.Hex())
		})
	})
}
