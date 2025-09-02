package etherscan

import "fmt"

type HttpStatusCodeError struct {
	URL        string
	StatusCode int
	Body       string
}

func (e *HttpStatusCodeError) Error() string {
	return fmt.Sprintf("HTTP status %s error: %d, body: %s", e.URL, e.StatusCode, e.Body)
}

type NetworkRequestError struct {
	Err error
}

func (e *NetworkRequestError) Error() string {
	return fmt.Sprintf("network request error: %v", e.Err)
}

type ContractVerificationInvalidStatusCodeError struct {
	URL        string
	StatusCode int
	Body       string
}

func (e *ContractVerificationInvalidStatusCodeError) Error() string {
	return fmt.Sprintf("invalid status code for %s: %d, body: %s", e.URL, e.StatusCode, e.Body)
}

type ContractVerificationMissingBytecodeError struct {
	APIURL          string
	ContractAddress string
}

func (e *ContractVerificationMissingBytecodeError) Error() string {
	return fmt.Sprintf("bytecode missing for contract %s on %s", e.ContractAddress, e.APIURL)
}

type ContractAlreadyVerifiedError struct {
	ContractName    string
	ContractAddress string
}

func (e *ContractAlreadyVerifiedError) Error() string {
	return fmt.Sprintf("contract %s (%s) already verified", e.ContractName, e.ContractAddress)
}

type VerifyError struct {
	Message string
}

func (e *VerifyError) Error() string {
	return fmt.Sprintf("verify error: %s", e.Message)
}
