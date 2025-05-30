package contractcall

import (
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
