package ethtype

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/donutnomad/eths/ecommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Legacy transaction (type 0x0)
const legacyTxJSON = `{
	"type": "0x0",
	"chainId": "0x1",
	"nonce": "0x15",
	"to": "0xdac17f958d2ee523a2206206994597c13d831ec7",
	"gas": "0x5208",
	"gasPrice": "0x4a817c800",
	"value": "0xde0b6b3a7640000",
	"input": "0x",
	"v": "0x25",
	"r": "0x1b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea",
	"s": "0x4ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c",
	"hash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
}`

// EIP-1559 transaction (type 0x2)
const eip1559TxJSON = `{
	"type": "0x2",
	"chainId": "0x1",
	"nonce": "0x42",
	"to": "0xdac17f958d2ee523a2206206994597c13d831ec7",
	"gas": "0x5208",
	"maxPriorityFeePerGas": "0x3b9aca00",
	"maxFeePerGas": "0x77359400",
	"value": "0x0",
	"input": "0xa9059cbb000000000000000000000000",
	"accessList": [],
	"v": "0x1",
	"r": "0x2b5e176d927f8e9ab405058b2d2457392da3e20f328b16ddabcebc33eaac5fea",
	"s": "0x3ba69724e8f69de52f0125ad8b3c5c2cef33019bac3249e2c0a2192766d1721c",
	"yParity": "0x1",
	"hash": "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
}`

// EIP-4844 blob transaction (type 0x3)
const blobTxJSON = `{
	"type": "0x3",
	"chainId": "0x1",
	"nonce": "0xa",
	"to": "0xdac17f958d2ee523a2206206994597c13d831ec7",
	"gas": "0x5208",
	"maxPriorityFeePerGas": "0x3b9aca00",
	"maxFeePerGas": "0x77359400",
	"maxFeePerBlobGas": "0x1",
	"value": "0x0",
	"input": "0x",
	"accessList": [],
	"blobVersionedHashes": ["0x0100000000000000000000000000000000000000000000000000000000000001"],
	"v": "0x0",
	"r": "0x1111111111111111111111111111111111111111111111111111111111111111",
	"s": "0x2222222222222222222222222222222222222222222222222222222222222222",
	"yParity": "0x0",
	"hash": "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
}`

func TestTx_UnmarshalJSON_Legacy(t *testing.T) {
	var tx Tx
	require.NoError(t, json.Unmarshal([]byte(legacyTxJSON), &tx))

	assert.Equal(t, uint8(0), tx.Type)
	assert.Equal(t, big.NewInt(1), tx.ChainID)
	assert.Equal(t, uint64(0x15), tx.Nonce)
	assert.Equal(t, ecommon.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"), *tx.To)
	assert.Equal(t, uint64(0x5208), tx.Gas)
	assert.Equal(t, big.NewInt(0x4a817c800), tx.GasPrice)
	assert.Equal(t, big.NewInt(0xde0b6b3a7640000), tx.Value)
	assert.Equal(t, []byte{}, tx.Input)
	assert.Equal(t, ecommon.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), tx.Hash)
}

func TestTx_UnmarshalJSON_EIP1559(t *testing.T) {
	var tx Tx
	require.NoError(t, json.Unmarshal([]byte(eip1559TxJSON), &tx))

	assert.Equal(t, uint8(2), tx.Type)
	assert.Equal(t, big.NewInt(1), tx.ChainID)
	assert.Equal(t, uint64(0x42), tx.Nonce)
	assert.Equal(t, uint64(0x5208), tx.Gas)
	assert.Equal(t, big.NewInt(0x3b9aca00), tx.MaxPriorityFeePerGas)
	assert.Equal(t, big.NewInt(0x77359400), tx.MaxFeePerGas)
	assert.True(t, tx.Value.Sign() == 0)
	assert.Equal(t, uint64(1), tx.YParity)
	//assert.NotNil(t, tx.AccessList)
	//assert.Empty(t, tx.AccessList)
}

func TestTx_UnmarshalJSON_Blob(t *testing.T) {
	var tx Tx
	require.NoError(t, json.Unmarshal([]byte(blobTxJSON), &tx))

	assert.Equal(t, uint8(3), tx.Type)
	assert.Equal(t, big.NewInt(1), tx.MaxFeePerBlobGas)
	assert.Len(t, tx.BlobVersionedHashes, 1)
	assert.Equal(t,
		ecommon.HexToHash("0x0100000000000000000000000000000000000000000000000000000000000001"),
		tx.BlobVersionedHashes[0],
	)
}

func TestTx_MarshalJSON_Roundtrip(t *testing.T) {
	tests := []struct {
		name    string
		rawJSON string
	}{
		{"legacy", legacyTxJSON},
		{"eip1559", eip1559TxJSON},
		{"blob", blobTxJSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unmarshal
			var tx Tx
			require.NoError(t, json.Unmarshal([]byte(tt.rawJSON), &tx))

			// Marshal
			encoded, err := json.Marshal(&tx)
			require.NoError(t, err)

			// Unmarshal again
			var tx2 Tx
			require.NoError(t, json.Unmarshal(encoded, &tx2))

			// Marshal again and compare JSON output for exact equivalence
			encoded2, err := json.Marshal(&tx2)
			require.NoError(t, err)
			assert.JSONEq(t, string(encoded), string(encoded2))
		})
	}
}

func TestTx_MarshalJSON_HexEncoding(t *testing.T) {
	var tx Tx
	require.NoError(t, json.Unmarshal([]byte(legacyTxJSON), &tx))

	encoded, err := json.Marshal(&tx)
	require.NoError(t, err)

	// Verify hex encoding in raw JSON
	var raw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(encoded, &raw))

	// type should be hex
	assert.Equal(t, `"0x0"`, string(raw["type"]))
	// nonce should be hex
	assert.Equal(t, `"0x15"`, string(raw["nonce"]))
	// gas should be hex
	assert.Equal(t, `"0x5208"`, string(raw["gas"]))
	// chainId should be hex
	assert.Equal(t, `"0x1"`, string(raw["chainId"]))
	// value should be hex
	assert.Equal(t, `"0xde0b6b3a7640000"`, string(raw["value"]))
	// input should be hex
	assert.Equal(t, `"0x"`, string(raw["input"]))
}

func TestTx_UnmarshalJSON_NilTo(t *testing.T) {
	// Contract creation: to is null
	raw := `{
		"type": "0x0",
		"nonce": "0x0",
		"to": null,
		"gas": "0x5208",
		"gasPrice": "0x1",
		"value": "0x0",
		"input": "0x6060",
		"v": "0x1b",
		"r": "0x1",
		"s": "0x2",
		"hash": "0x0000000000000000000000000000000000000000000000000000000000000001"
	}`
	var tx Tx
	require.NoError(t, json.Unmarshal([]byte(raw), &tx))
	assert.Nil(t, tx.To)
	assert.Equal(t, []byte{0x60, 0x60}, tx.Input)
}

func TestTxDetail_UnmarshalJSON(t *testing.T) {
	raw := `{
		"type": "0x2",
		"chainId": "0x1",
		"nonce": "0x42",
		"to": "0xdac17f958d2ee523a2206206994597c13d831ec7",
		"gas": "0x5208",
		"maxPriorityFeePerGas": "0x3b9aca00",
		"maxFeePerGas": "0x77359400",
		"value": "0x0",
		"input": "0xa9059cbb",
		"v": "0x1",
		"r": "0x1",
		"s": "0x2",
		"yParity": "0x1",
		"hash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"root": "0x",
		"status": "0x1",
		"cumulativeGasUsed": "0xb41e",
		"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"logs": [],
		"contractAddress": "0x0000000000000000000000000000000000000000",
		"gasUsed": "0x5208",
		"effectiveGasPrice": "0x3b9aca00",
		"blockHash": "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		"blockNumber": "0x100",
		"transactionIndex": "0x3"
	}`

	var detail TxDetail
	require.NoError(t, json.Unmarshal([]byte(raw), &detail))

	// Tx fields
	assert.Equal(t, uint8(2), detail.Tx.Type)
	assert.Equal(t, big.NewInt(1), detail.Tx.ChainID)
	assert.Equal(t, uint64(0x42), detail.Tx.Nonce)
	assert.Equal(t, ecommon.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7"), *detail.Tx.To)
	assert.Equal(t, uint64(0x5208), detail.Tx.Gas)
	assert.Equal(t, big.NewInt(0x3b9aca00), detail.Tx.MaxPriorityFeePerGas)
	assert.Equal(t, big.NewInt(0x77359400), detail.Tx.MaxFeePerGas)
	assert.Equal(t, []byte{0xa9, 0x05, 0x9c, 0xbb}, detail.Tx.Input)
	assert.Equal(t, ecommon.HexToHash("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), detail.Tx.Hash)

	// Receipt fields
	assert.Equal(t, uint64(1), detail.Status)
	assert.Equal(t, uint64(0xb41e), detail.CumulativeGasUsed)
	assert.Equal(t, uint64(0x5208), detail.GasUsed)
	assert.Equal(t, big.NewInt(0x3b9aca00), detail.EffectiveGasPrice)
	assert.Equal(t, big.NewInt(0x100), detail.BlockNumber)
	assert.Equal(t, uint(3), detail.TransactionIndex)
	assert.Equal(t, ecommon.HexToHash("0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"), detail.BlockHash)
}

func TestTxDetail_MarshalJSON_Roundtrip(t *testing.T) {
	raw := `{
		"type": "0x2",
		"chainId": "0x1",
		"nonce": "0x42",
		"to": "0xdac17f958d2ee523a2206206994597c13d831ec7",
		"gas": "0x5208",
		"maxPriorityFeePerGas": "0x3b9aca00",
		"maxFeePerGas": "0x77359400",
		"value": "0x0",
		"input": "0xa9059cbb",
		"v": "0x1",
		"r": "0x1",
		"s": "0x2",
		"yParity": "0x1",
		"hash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"root": "0x",
		"status": "0x1",
		"cumulativeGasUsed": "0xb41e",
		"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"logs": [],
		"contractAddress": "0x0000000000000000000000000000000000000000",
		"gasUsed": "0x5208",
		"effectiveGasPrice": "0x3b9aca00",
		"blockHash": "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		"blockNumber": "0x100",
		"transactionIndex": "0x3"
	}`

	var detail TxDetail
	require.NoError(t, json.Unmarshal([]byte(raw), &detail))

	encoded, err := json.Marshal(&detail)
	require.NoError(t, err)

	var detail2 TxDetail
	require.NoError(t, json.Unmarshal(encoded, &detail2))

	encoded2, err := json.Marshal(&detail2)
	require.NoError(t, err)
	assert.JSONEq(t, string(encoded), string(encoded2))
}

func TestTxDetail_MarshalJSON_HexEncoding(t *testing.T) {
	raw := `{
		"type": "0x2",
		"chainId": "0x1",
		"nonce": "0x42",
		"to": "0xdac17f958d2ee523a2206206994597c13d831ec7",
		"gas": "0x5208",
		"maxPriorityFeePerGas": "0x3b9aca00",
		"maxFeePerGas": "0x77359400",
		"value": "0x0",
		"input": "0x",
		"v": "0x1",
		"r": "0x1",
		"s": "0x2",
		"hash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"root": "0x",
		"status": "0x1",
		"cumulativeGasUsed": "0xb41e",
		"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"logs": [],
		"contractAddress": "0x0000000000000000000000000000000000000000",
		"gasUsed": "0x5208",
		"effectiveGasPrice": "0x3b9aca00",
		"blockNumber": "0x100",
		"transactionIndex": "0x3"
	}`

	var detail TxDetail
	require.NoError(t, json.Unmarshal([]byte(raw), &detail))

	encoded, err := json.Marshal(&detail)
	require.NoError(t, err)

	var fields map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(encoded, &fields))

	// Tx fields are hex
	assert.Equal(t, `"0x2"`, string(fields["type"]))
	assert.Equal(t, `"0x42"`, string(fields["nonce"]))
	assert.Equal(t, `"0x5208"`, string(fields["gas"]))
	// Receipt fields are hex
	assert.Equal(t, `"0x1"`, string(fields["status"]))
	assert.Equal(t, `"0xb41e"`, string(fields["cumulativeGasUsed"]))
	assert.Equal(t, `"0x5208"`, string(fields["gasUsed"]))
	assert.Equal(t, `"0x3b9aca00"`, string(fields["effectiveGasPrice"]))
	assert.Equal(t, `"0x100"`, string(fields["blockNumber"]))
	assert.Equal(t, `"0x3"`, string(fields["transactionIndex"]))
}

func TestTxReceipt_MarshalJSON_Roundtrip(t *testing.T) {
	raw := `{
		"type": "0x2",
		"root": "0x",
		"status": "0x1",
		"cumulativeGasUsed": "0x5208",
		"logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"logs": [],
		"transactionHash": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"contractAddress": "0x0000000000000000000000000000000000000000",
		"gasUsed": "0x5208",
		"effectiveGasPrice": "0x3b9aca00",
		"blockHash": "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		"blockNumber": "0x100",
		"transactionIndex": "0x0"
	}`

	var receipt TxReceipt
	require.NoError(t, json.Unmarshal([]byte(raw), &receipt))

	assert.Equal(t, uint8(2), receipt.Type)
	assert.Equal(t, uint64(1), receipt.Status)
	assert.Equal(t, uint64(0x5208), receipt.CumulativeGasUsed)
	assert.Equal(t, uint64(0x5208), receipt.GasUsed)
	assert.Equal(t, big.NewInt(0x3b9aca00), receipt.EffectiveGasPrice)
	assert.Equal(t, big.NewInt(0x100), receipt.BlockNumber)

	// Roundtrip
	encoded, err := json.Marshal(&receipt)
	require.NoError(t, err)

	var receipt2 TxReceipt
	require.NoError(t, json.Unmarshal(encoded, &receipt2))
	assert.Equal(t, receipt, receipt2)
}
