package etherscan

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type EtherscanClient struct {
	client  *resty.Client
	apiKey  string
	baseURL string
}

func NewEtherscanClient(apiKey string) *EtherscanClient {
	return &EtherscanClient{
		client:  resty.New(),
		apiKey:  apiKey,
		baseURL: "https://api.etherscan.io/v2/api",
	}
}

func NewEtherscanClientWith(baseUrl, apiKey string, client *resty.Client) *EtherscanClient {
	return &EtherscanClient{
		client:  client,
		apiKey:  apiKey,
		baseURL: baseUrl,
	}
}

func (e *EtherscanClient) SetDebug() {
	e.client.SetDebug(true)
}

func (e *EtherscanClient) SetClient(client *resty.Client) {
	e.client = client
}

func getResult[T any](ctx context.Context,
	client *EtherscanClient,
	params map[string]string,
	module string,
	action string,
	chainID uint64) (*Response[T], error) {
	var response Response[T]
	resp, err := client.client.R().
		SetContext(ctx).
		SetQueryParam("apikey", client.apiKey).
		SetQueryParam("chainid", strconv.FormatUint(chainID, 10)).
		SetQueryParam("module", module).
		SetQueryParam("action", action).
		SetQueryParams(params).
		SetResult(&response).
		Get(client.baseURL)
	if err != nil {
		return nil, &NetworkRequestError{Err: err}
	}
	if resp.StatusCode() != 200 {
		return nil, &HttpStatusCodeError{
			URL:        resp.Request.URL,
			StatusCode: resp.StatusCode(),
			Body:       string(resp.Body()),
		}
	}
	if !response.IsSuccess() && !strings.Contains(response.Message, "No records found") {
		// example: etherscan api returned error: https://api.etherscan.io/v2/api?action=getLogs&apikey=xxx&chainid=1&module=logs&offset=10&page=1 NOTOK
		return nil, fmt.Errorf("etherscan api returned error: %s %s", resp.Request.URL, response.Message)
	}
	return &response, nil
}
