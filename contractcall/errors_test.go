package contractcall

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/donutnomad/blockchain-alg/xecdsa"
	"github.com/donutnomad/eths/common"
	"github.com/donutnomad/eths/common/hexutil"
	"github.com/donutnomad/eths/ecommon"
	"github.com/donutnomad/eths/ethclient"
	"github.com/donutnomad/eths/ethtype"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/pkg/errors"
	"testing"
)

func TestErrEstimateGasError(t *testing.T) {
	err := providerErr1()
	var e1 *EstimateGasError
	if errors.As(err, &e1) {
		t.Log("ok")
	} else {
		t.Fatal("failed")
	}
}

func providerErr1() error {
	return &EstimateGasError{Err: errors.New("err1")}
}
