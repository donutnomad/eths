package etherscan

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func (e *EtherscanClient) VerifyAndCheck(ctx context.Context, chainID uint64,
	input json.RawMessage,
	solcVersion string,
	verifyContractName string,
	contractAddress common.Address,
	constructorAbiEncode []byte,
) error {
	response, err := e.Verify(ctx, chainID, contractAddress, string(input), verifyContractName, solcVersion, constructorAbiEncode)
	if err != nil {
		var alreadyVerifiedErr *ContractAlreadyVerifiedError
		if errors.As(err, &alreadyVerifiedErr) {
			return nil
		}
		return err
	}
	guid := response.Result
	var ticker = time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return ctx.Err()
		}
		status, err := e.CheckVerifyStatus(ctx, chainID, guid)
		if err != nil {
			return fmt.Errorf("check verify status %w, guid: %s", err, guid)
		}
		if status.IsPending() {
			continue
		}
		if status.IsAlreadyVerified() {
			return nil
		}
		if status.IsFailure() {
			return fmt.Errorf("fail - Unable to verify %w, guid: %s", err, guid)
		}
		return nil
	}
}

func (e *EtherscanClient) Verify(
	ctx context.Context,
	chainID uint64,
	contractAddress common.Address,
	sourceCode string,
	contractName string,
	compilerVersion string,
	constructorAbiEncode []byte,
) (*VerifyResponse, error) {
	var verifyResp VerifyResponse
	resp, err := e.client.R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetQueryParam("chainid", strconv.FormatUint(chainID, 10)).
		SetFormDataFromValues(url.Values{
			"apikey":                []string{e.apiKey},
			"module":                []string{"contract"},
			"action":                []string{"verifysourcecode"},
			"contractaddress":       []string{contractAddress.Hex()},
			"sourceCode":            []string{sourceCode},
			"codeformat":            []string{"solidity-standard-json-input"},
			"contractname":          []string{contractName},
			"compilerversion":       []string{compilerVersion},
			"constructorArguements": []string{hex.EncodeToString(constructorAbiEncode)},
		}).
		SetResult(&verifyResp).
		Post(e.baseURL)
	if err != nil {
		return nil, &NetworkRequestError{Err: err}
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, &ContractVerificationInvalidStatusCodeError{
			URL:        resp.Request.URL,
			StatusCode: resp.StatusCode(),
			Body:       string(resp.Body()),
		}
	}

	if verifyResp.IsBytecodeMissingInNetworkError() {
		return nil, &ContractVerificationMissingBytecodeError{
			APIURL:          e.baseURL,
			ContractAddress: contractAddress.Hex(),
		}
	} else if verifyResp.IsAlreadyVerified() {
		return nil, &ContractAlreadyVerifiedError{
			ContractName:    contractName,
			ContractAddress: contractAddress.Hex(),
		}
	} else if !verifyResp.IsSuccess() {
		return nil, &VerifyError{
			Message: verifyResp.Message,
		}
	}
	return &verifyResp, nil
}

func (e *EtherscanClient) CheckVerifyStatus(ctx context.Context, chainID uint64, guid string) (*VerifyResponse, error) {
	resp, err := getResult[string](
		ctx,
		e,
		map[string]string{
			"guid": guid,
		},
		"contract",
		"checkverifystatus",
		chainID,
	)
	if err != nil {
		return nil, err
	}
	return &VerifyResponse{
		Response: *resp,
	}, nil
}

type VerifyResponse struct {
	Response[string]
}

func (e *VerifyResponse) IsPending() bool {
	return e.Message == "Pending in queue"
}

func (e *VerifyResponse) IsFailure() bool {
	return e.Message == "Fail - Unable to verify"
}

func (e *VerifyResponse) IsVerified() bool {
	return e.Message == "Pass - Verified"
}

func (e *VerifyResponse) IsBytecodeMissingInNetworkError() bool {
	return strings.HasPrefix(e.Message, "Unable to locate ContractCode at")
}

func (e *VerifyResponse) IsAlreadyVerified() bool {
	return strings.HasPrefix(e.Message, "Smart-contract already verified") ||
		strings.HasPrefix(e.Message, "Contract source code already verified") ||
		strings.HasPrefix(e.Message, "Already Verified")
}
