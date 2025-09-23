package ethtype

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

type AddressRequest struct {
	Address Address `json:"address" form:"address"`
}

func TestAddressJson(t *testing.T) {
	makeRequest(map[string]string{
		"address": "0x1234567890abcdef1234567890abcdef12345678",
	}, func(req AddressRequest) {
		assert.Equal(t, "0x1234567890AbcdEF1234567890aBcdef12345678", req.Address.Hex())
	})
}

func TestAddressForm(t *testing.T) {
	makeFormRequest(map[string]string{
		"address": "0x1234567890abcdef1234567890abcdef12345678",
	}, func(req AddressRequest) {
		assert.Equal(t, "0x1234567890AbcdEF1234567890aBcdef12345678", req.Address.Hex())
	})
}

func TestAddressFormFailed(t *testing.T) {
	assert.Panics(t, func() {
		makeFormRequest(map[string]string{
			"address": "0x1234567890abcdef1234567890abcdef1234567878",
		}, func(req AddressRequest) {
			spew.Dump(req)
		})
	})
}

// 正常
func TestAddressQuery(t *testing.T) {
	makeGetRequest(map[string]string{
		"address": "0x1234567890abcdef1234567890abcdef12345678",
	}, func(req AddressRequest) {
		spew.Dump(req)
		assert.Equal(t, "0x1234567890AbcdEF1234567890aBcdef12345678", req.Address.Hex())
	})
}

// 超过32个字节，报错
func TestAddressQuery2(t *testing.T) {
	assert.Panics(t, func() {
		makeGetRequest(map[string]string{
			"address": "0x1234567890abcdef1234567890abcdef1234567878",
		}, func(req AddressRequest) {
			assert.Equal(t, "0x1234567890AbcdEF1234567890aBcdef12345678", req.Address.Hex())
		})
	})
}
