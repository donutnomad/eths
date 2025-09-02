package etherscan

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

func (r *Response[T]) IsSuccess() bool {
	return r.Status == "1"
}

func (r *Response[T]) UnmarshalJSON(data []byte) error {
	type Tmp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	var tmp Tmp
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	r.Status = tmp.Status
	r.Message = tmp.Message
	if tmp.Status == "1" {
		type Tmp2 struct {
			Result T `json:"result"`
		}
		var tmp2 Tmp2
		if err := json.Unmarshal(data, &tmp2); err != nil {
			return err
		}
		r.Result = tmp2.Result
	}
	return nil
}

type Uint64 uint64

func (u *Uint64) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if hexStr == "" {
		*u = 0
		return nil
	}
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	parseUint, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return fmt.Errorf("etherscan uint64 %s %w", hexStr, err)
	}
	*u = Uint64(parseUint)
	return nil
}

func (u Uint64) MarshalJSON() ([]byte, error) {
	return json.Marshal("0x" + strconv.FormatUint(uint64(u), 16))
}

type HexBs []byte

func (h *HexBs) UnmarshalJSON(data []byte) error {
	var hexStr string
	if err := json.Unmarshal(data, &hexStr); err != nil {
		return err
	}
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	bs, err := hex.DecodeString(hexStr)
	if err != nil {
		return err
	}
	*h = bs
	return nil
}

func (h HexBs) MarshalJSON() ([]byte, error) {
	return json.Marshal("0x" + hex.EncodeToString(h))
}
