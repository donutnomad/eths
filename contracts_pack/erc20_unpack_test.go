package contracts_pack

import (
	"encoding/hex"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestERC20UnpackTransferInput(t *testing.T) {
	data, err := hex.DecodeString("a9059cbb000000000000000000000000e5b877b55855e009ac689dce4f435304abd241120000000000000000000000000000000000000000000000000000000000002710")
	if err != nil {
		panic(err)
	}
	eRC20 := NewERC20()
	to, value, err := eRC20.UnpackInputTransfer(data)
	if err != nil {
		panic(err)
	}
	spew.Dump(to, value)
}
