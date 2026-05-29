package gogenapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/donutnomad/eths/abigen/abigen2/abigen"
	ethabi "github.com/ethereum/go-ethereum/accounts/abi"
)

type InputKind int

const (
	InputKindABI InputKind = iota + 1
	InputKindArtifact
)

type Input struct {
	ABI          []byte
	Bytecode     []byte
	ContractName string
	Kind         InputKind
	Path         string
}

type GenerateOptions struct {
	ABIPath     string
	PackageName string
	TypeName    string
	Aliases     map[string]string
}

func LoadInput(path string) (*Input, error) {
	path = strings.TrimPrefix(path, "@")
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("empty ABI input: %s", path)
	}
	if trimmed[0] != '{' {
		return &Input{
			ABI:  raw,
			Kind: InputKindABI,
			Path: path,
		}, nil
	}

	var artifact struct {
		ContractName string          `json:"contractName"`
		ABI          json.RawMessage `json:"abi"`
		Bytecode     string          `json:"bytecode"`
	}
	if err := json.Unmarshal(raw, &artifact); err != nil {
		return nil, err
	}
	abiRaw := bytes.TrimSpace(artifact.ABI)
	if len(abiRaw) == 0 || bytes.Equal(abiRaw, []byte("null")) {
		return nil, fmt.Errorf("artifact %s does not contain abi", path)
	}
	return &Input{
		ABI:          artifact.ABI,
		Bytecode:     []byte(artifact.Bytecode),
		ContractName: artifact.ContractName,
		Kind:         InputKindArtifact,
		Path:         path,
	}, nil
}

func DefaultTypeName(input *Input) string {
	if input == nil {
		return ""
	}
	if input.Kind == InputKindArtifact && input.ContractName != "" {
		return ethabi.ToCamelCase(input.ContractName)
	}
	base := filepath.Base(input.Path)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return ethabi.ToCamelCase(name)
}

func Generate(opts GenerateOptions) ([]byte, error) {
	if opts.ABIPath == "" {
		return nil, fmt.Errorf("ABIPath is required")
	}
	if opts.PackageName == "" {
		return nil, fmt.Errorf("PackageName is required")
	}
	input, err := LoadInput(opts.ABIPath)
	if err != nil {
		return nil, err
	}
	typeName := opts.TypeName
	if typeName == "" {
		typeName = DefaultTypeName(input)
	}
	if typeName == "" {
		return nil, fmt.Errorf("TypeName is required")
	}
	aliases := opts.Aliases
	if aliases == nil {
		aliases = map[string]string{}
	}
	code, err := abigen.BindV2(
		[]string{typeName},
		[]string{string(input.ABI)},
		[]string{string(input.Bytecode)},
		opts.PackageName,
		map[string]string{},
		aliases,
	)
	if err != nil {
		return nil, err
	}
	return []byte(code), nil
}
