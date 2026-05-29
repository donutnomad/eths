// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.Package}}

import (
	"bytes"
	"context"
	"errors"
	"math/big"

	"github.com/donutnomad/eths/contractcall"
	"github.com/donutnomad/eths/abi"
	"github.com/donutnomad/eths/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = bytes.Equal
	_ = context.Background
	_ = errors.New
	_ = big.NewInt
	_ = contractcall.SendTxE
	_ = common.Big1
	_ = types.BloomLookup
	_ = abi.ConvertType
)

var (
	errNoEventSignature       = errors.New("missing event signature")
	errEventSignatureMismatch = errors.New("event signature mismatch")
	errMethodSignatureMismatch = errors.New("method signature mismatch")
	errErrorSignatureMismatch  = errors.New("error signature mismatch")
)

type ContractEvent interface {
	Name() string
	Topic0() common.Hash
}

type ContractError interface {
	Name() string
	ErrorID() common.Hash
}

{{$structs := .Structs}}
{{range $structs}}
	// {{.Name}} is an auto generated low-level Go binding around an user-defined struct.
	type {{.Name}} struct {
	{{range $field := .Fields}}
	{{capitalise $field.Name}} {{$field.Type}}{{end}}
	}
{{end}}

{{range $contract := .Contracts}}
	// {{.Type}}MetaData contains all meta data concerning the {{.Type}} contract.
	var {{.Type}}MetaData = bind.MetaData{
		ABI: "{{.InputABI}}",
		{{if (index $.Libraries .Type) -}}
		ID: "{{index $.Libraries .Type}}",
		{{ else -}}
		ID: "{{.Type}}",
		{{end -}}
		{{if .InputBin -}}
		Bin: "0x{{.InputBin}}",
		{{end -}}
		{{if .Libraries -}}
		Deps: []*bind.MetaData{
		{{- range $name, $pattern := .Libraries}}
			&{{$name}}MetaData,
		{{- end}}
		},
		{{end}}
	}

	var _{{.Type}}ABI *abi.ABI

	func init() {
		var err error
		_{{.Type}}ABI, err = {{.Type}}MetaData.ParseABI()
		if err != nil {
			panic(errors.New("invalid ABI: " + err.Error()))
		}
	}

	// {{.Type}} is an auto generated Go binding around an Ethereum contract.
	type {{.Type}} struct {
		{{range .Calls}}
		{{.Normalized.Name}} {{methodtypename $contract.Type .Normalized.Name}}
		{{- end}}
		{{range .Events}}
		{{.Normalized.Name}}Event {{eventtypename $contract.Type .Normalized.Name}}
		{{- end}}
		{{range .Errors}}
		{{.Normalized.Name}}Error {{errortypename $contract.Type .Normalized.Name}}
		{{- end}}
	}

	// New{{.Type}} creates a new instance of {{.Type}}.
	func New{{.Type}}() *{{.Type}} {
		return &{{.Type}}{}
	}

	// Instance creates a wrapper for a deployed contract instance at the given address.
	// Use this to create the instance object passed to abigen v2 library functions Call, Transact, etc.
	func (c *{{.Type}}) Instance(backend bind.ContractBackend, addr common.Address) *bind.BoundContract {
		 return bind.NewBoundContract(addr, *_{{.Type}}ABI, backend, backend, backend)
	}

	{{ if .Constructor.Inputs }}
	// PackConstructor is the Go binding used to pack the parameters required for
	// contract deployment.
	// Solidity: {{.Constructor.String}}
	func ({{ decapitalise $contract.Type}} *{{$contract.Type}}) PackConstructor({{range .Constructor.Inputs}} {{.Name}} {{bindtype .Type $structs}}, {{end}}) []byte {
		return {{decapitalise $contract.Type}}Must(_{{$contract.Type}}ABI.Pack("" {{range .Constructor.Inputs}}, {{.Name}}{{end}}))
	}
	{{ end }}

	{{range .Calls}}
		// {{methodtypename $contract.Type .Normalized.Name}} is the Go binding around the {{.Original.Name}} contract method.
		// MethodID: 0x{{printf "%x" .Original.ID}}
	// Solidity: {{.Original.String}}
	type {{methodtypename $contract.Type .Normalized.Name}} struct {
	}

		{{- if .Normalized.Inputs}}
			// {{inputtypename $contract.Type .Normalized.Name}} is the Go binding for the {{.Original.Name}} method input parameters.
			type {{inputtypename $contract.Type .Normalized.Name}} struct {
				{{- range .Normalized.Inputs}}
				{{inputfieldname .Name}} {{bindtype .Type $structs}}
				{{- end}}
			}
		{{- end}}

		func (m {{methodtypename $contract.Type .Normalized.Name}}) Name() string { return "{{.Original.Name}}" }
		func (m {{methodtypename $contract.Type .Normalized.Name}}) MethodID() [4]byte { return {{methodidbytes .Original.ID}} }
		{{- if useinputstruct .}}
		func (m {{methodtypename $contract.Type .Normalized.Name}}) Pack(input {{inputtypename $contract.Type .Normalized.Name}}) []byte {
			return {{decapitalise $contract.Type}}Must(m.TryPack(input))
		{{- else}}
		func (m {{methodtypename $contract.Type .Normalized.Name}}) Pack({{range .Normalized.Inputs}} {{.Name}} {{bindtype .Type $structs}}, {{end}}) []byte {
			return {{decapitalise $contract.Type}}Must(m.TryPack({{range .Normalized.Inputs}} {{.Name}}, {{end}}))
		{{- end}}
		}
	{{- if useinputstruct .}}
	func (m {{methodtypename $contract.Type .Normalized.Name}}) TryPack(input {{inputtypename $contract.Type .Normalized.Name}}) ([]byte, error) {
		return _{{$contract.Type}}ABI.Pack(m.Name() {{range .Normalized.Inputs}}, input.{{inputfieldname .Name}}{{end}})
	}
	{{- else}}
	func (m {{methodtypename $contract.Type .Normalized.Name}}) TryPack({{range .Normalized.Inputs}} {{.Name}} {{bindtype .Type $structs}}, {{end}}) ([]byte, error) {
		return _{{$contract.Type}}ABI.Pack(m.Name() {{range .Normalized.Inputs}}, {{.Name}}{{end}})
	}
		{{- end}}
	{{- if useinputstruct .}}
	func (m {{methodtypename $contract.Type .Normalized.Name}}) UnpackInput(callData []byte) (*{{inputtypename $contract.Type .Normalized.Name}}, error) {
		return {{decapitalise $contract.Type}}UnpackInput[{{inputtypename $contract.Type .Normalized.Name}}](m.Name(), callData)
	}
		{{- else if .Normalized.Inputs}}
		func (m {{methodtypename $contract.Type .Normalized.Name}}) UnpackInput(callData []byte) ({{range .Normalized.Inputs}} {{.Name}} {{bindtype .Type $structs}}, {{end}} err error) {
			input, err := {{decapitalise $contract.Type}}UnpackInput[{{inputtypename $contract.Type .Normalized.Name}}](m.Name(), callData)
			if err != nil {
				return {{range .Normalized.Inputs}}{{if ispointertype .Type}}new({{underlyingbindtype .Type }}), {{else}}*new({{bindtype .Type $structs}}), {{end}}{{end}} err
			}
			return {{range .Normalized.Inputs}}input.{{inputfieldname .Name}}, {{end}} nil
			}
	{{- end}}
	{{- if methodisconstant .}}
		{{- if .Normalized.Outputs}}
			{{- if useinputstruct .}}
			func (m {{methodtypename $contract.Type .Normalized.Name}}) Call(ctx context.Context, client bind.ContractCaller, addr common.Address, opts *bind.CallOpts, input {{inputtypename $contract.Type .Normalized.Name}}) (
			{{- if .Structured}} {{outputtypename $contract.Type .Normalized.Name}},{{else}}
			{{- range .Normalized.Outputs}} {{bindtype .Type $structs}},{{- end }}
			{{- end }} error) {
				output, err := {{decapitalise $contract.Type}}CallRaw(ctx, client, addr, opts, m.Pack(input))
				if err != nil {
					{{- if .Structured}}
					return *new({{outputtypename $contract.Type .Normalized.Name}}), err
					{{- else}}
					return {{range .Normalized.Outputs}}{{if ispointertype .Type}}new({{underlyingbindtype .Type }}), {{else}}*new({{bindtype .Type $structs}}), {{end}}{{end}} err
					{{- end}}
				}
				return m.Unpack(output)
			}
			{{- else}}
			func (m {{methodtypename $contract.Type .Normalized.Name}}) Call(ctx context.Context, client bind.ContractCaller, addr common.Address, opts *bind.CallOpts, {{range .Normalized.Inputs}} {{.Name}} {{bindtype .Type $structs}}, {{end}}) (
			{{- if .Structured}} {{outputtypename $contract.Type .Normalized.Name}},{{else}}
			{{- range .Normalized.Outputs}} {{bindtype .Type $structs}},{{- end }}
			{{- end }} error) {
				output, err := {{decapitalise $contract.Type}}CallRaw(ctx, client, addr, opts, m.Pack({{range .Normalized.Inputs}} {{.Name}}, {{end}}))
				if err != nil {
					{{- if .Structured}}
					return *new({{outputtypename $contract.Type .Normalized.Name}}), err
					{{- else}}
					return {{range .Normalized.Outputs}}{{if ispointertype .Type}}new({{underlyingbindtype .Type }}), {{else}}*new({{bindtype .Type $structs}}), {{end}}{{end}} err
					{{- end}}
				}
				return m.Unpack(output)
			}
			{{- end}}
		{{- else}}
			{{- if useinputstruct .}}
			func (m {{methodtypename $contract.Type .Normalized.Name}}) Call(ctx context.Context, client bind.ContractCaller, addr common.Address, opts *bind.CallOpts, input {{inputtypename $contract.Type .Normalized.Name}}) error {
				_, err := {{decapitalise $contract.Type}}CallRaw(ctx, client, addr, opts, m.Pack(input))
				return err
			}
			{{- else}}
			func (m {{methodtypename $contract.Type .Normalized.Name}}) Call(ctx context.Context, client bind.ContractCaller, addr common.Address, opts *bind.CallOpts, {{range .Normalized.Inputs}} {{.Name}} {{bindtype .Type $structs}}, {{end}}) error {
				_, err := {{decapitalise $contract.Type}}CallRaw(ctx, client, addr, opts, m.Pack({{range .Normalized.Inputs}} {{.Name}}, {{end}}))
				return err
			}
			{{- end}}
		{{- end}}
	{{- else}}
		{{- if useinputstruct .}}
		func (m {{methodtypename $contract.Type .Normalized.Name}}) Call(
			ctx context.Context,
			client contractcall.ISendTxClient,
			addr common.Address,
			tx contractcall.TxOpts,
			input {{inputtypename $contract.Type .Normalized.Name}},
		) (*types.Transaction, error) {
			return contractcall.SendTxE(ctx, client, tx.ChainID, tx.Value, m.Pack(input), &addr, tx.Payer, tx.Manager, tx.BeforeSend, tx.NoSend, true, tx.Options...)
		}
		{{- else}}
		func (m {{methodtypename $contract.Type .Normalized.Name}}) Call(
			ctx context.Context,
			client contractcall.ISendTxClient,
			addr common.Address,
			tx contractcall.TxOpts,
			{{- range .Normalized.Inputs}}
			{{.Name}} {{bindtype .Type $structs}},
			{{- end}}
		) (*types.Transaction, error) {
			return contractcall.SendTxE(ctx, client, tx.ChainID, tx.Value, m.Pack({{range .Normalized.Inputs}} {{.Name}}, {{end}}), &addr, tx.Payer, tx.Manager, tx.BeforeSend, tx.NoSend, true, tx.Options...)
		}
		{{- end}}
	{{- end}}
		{{/* Unpack method is needed only when there are return args */}}
		{{if .Normalized.Outputs }}
			{{ if .Structured }}
			// {{outputtypename $contract.Type .Normalized.Name}} serves as a container for the return parameters of contract
			// method {{ .Normalized.Name }}.
			type {{outputtypename $contract.Type .Normalized.Name}} struct {
			  {{range .Normalized.Outputs}}
			  {{capitalise .Name}} {{bindtype .Type $structs}}{{end}}
			}
			{{ end }}
			func (m {{methodtypename $contract.Type .Normalized.Name}}) Unpack(data []byte) (
		{{- if .Structured}} {{outputtypename $contract.Type .Normalized.Name}},{{else}}
		{{- range .Normalized.Outputs}} {{bindtype .Type $structs}},{{- end }}
		{{- end }} error) {
		out, err := _{{$contract.Type}}ABI.Unpack(m.Name(), data)
				{{- if .Structured}}
				outstruct := new({{outputtypename $contract.Type .Normalized.Name}})
					if err != nil {
						return *outstruct, err
					}
					{{- range $i, $t := .Normalized.Outputs}}
					outstruct.{{capitalise .Name}} = {{decapitalise $contract.Type}}ConvertType[{{bindtype .Type $structs}}](out[{{$i}}])
					{{- end }}
					return *outstruct, nil{{else}}
					if err != nil {
						return {{range $i, $_ := .Normalized.Outputs}}{{if ispointertype .Type}}new({{underlyingbindtype .Type }}), {{else}}*new({{bindtype .Type $structs}}), {{end}}{{end}} err
					}
					{{- range $i, $t := .Normalized.Outputs}}
					out{{$i}} := {{decapitalise $contract.Type}}ConvertType[{{bindtype .Type $structs}}](out[{{$i}}])
					{{- end}}
					return {{range $i, $t := .Normalized.Outputs}}out{{$i}}, {{end}} nil
				{{- end}}
			}
		{{end}}
	{{- if .Normalized.Inputs}}
			func (input *{{inputtypename $contract.Type .Normalized.Name}}) unpack(arguments []any) {
			{{- range $i, $t := .Normalized.Inputs}}
			input.{{inputfieldname .Name}} = {{decapitalise $contract.Type}}ConvertType[{{bindtype .Type $structs}}](arguments[{{$i}}])
			{{- end}}
			}
	{{- end}}
	{{end}}

	{{ if .Calls }}
	// UnpackInput attempts to decode the provided call data using contract method definitions.
	func ({{ decapitalise $contract.Type}} *{{$contract.Type}}) UnpackInput(callData []byte) (any, error) {
		if len(callData) < 4 {
			return nil, errMethodSignatureMismatch
		}
		{{- range .Calls}}
		if bytes.Equal(callData[:4], _{{$contract.Type}}ABI.Methods["{{.Original.Name}}"].ID[:4]) {
			{{- if .Normalized.Inputs}}
			return {{decapitalise $contract.Type}}UnpackInput[{{inputtypename $contract.Type .Normalized.Name}}]("{{.Original.Name}}", callData)
			{{- else}}
			return {{ decapitalise $contract.Type}}.{{.Normalized.Name}}, nil
			{{- end}}
		}
		{{- end}}
		return nil, errMethodSignatureMismatch
	}
	{{ end }}

	{{range .Events}}
		// {{eventtypename $contract.Type .Normalized.Name}} is the Go binding around the {{.Original.Name}} contract event.
		// Topic0: {{.Original.ID}}
	// Solidity: {{.Original.String}}
	type {{eventtypename $contract.Type .Normalized.Name}} struct {
		{{- range .Normalized.Inputs}}
			{{ capitalise .Name}}
			{{- if .Indexed}} {{ bindtopictype .Type $structs}}{{- else}} {{ bindtype .Type $structs}}{{ end }}
		{{- end}}
		Raw *types.Log // Blockchain specific contextual infos
	}

	const {{eventtypename $contract.Type .Normalized.Name}}Name = "{{.Original.Name}}"
	func (e {{eventtypename $contract.Type .Normalized.Name}}) Name() string { return {{eventtypename $contract.Type .Normalized.Name}}Name }
	func (e {{eventtypename $contract.Type .Normalized.Name}}) Topic0() common.Hash { return common.HexToHash("{{.Original.ID}}") }
	func (e {{eventtypename $contract.Type .Normalized.Name}}) Unpack(log *types.Log) (*{{eventtypename $contract.Type .Normalized.Name}}, error) {
		return {{decapitalise $contract.Type}}UnpackEvent[{{eventtypename $contract.Type .Normalized.Name}}](e.Name(), log)
	}
	func (out *{{eventtypename $contract.Type .Normalized.Name}}) unpack(values []any) {
		{{- $dataIndex := 0}}
		{{- range .Normalized.Inputs}}
		{{- if not .Indexed}}
		out.{{capitalise .Name}} = {{decapitalise $contract.Type}}ConvertType[{{bindtype .Type $structs}}](values[{{$dataIndex}}])
		{{- $dataIndex = add $dataIndex 1}}
		{{- end}}
		{{- end}}
	}
	func (out *{{eventtypename $contract.Type .Normalized.Name}}) setRaw(log *types.Log) { out.Raw = log }
	{{end}}

	{{range .Errors}}
		// {{errortypename $contract.Type .Normalized.Name}} is the Go binding around the {{.Original.Name}} contract error.
		// ErrorID: {{.Original.ID}}
	// Solidity: {{.Original.String}}
	type {{errortypename $contract.Type .Normalized.Name}} struct {
		{{- range .Normalized.Inputs}}
		{{capitalise .Name}} {{if .Indexed}}{{bindtopictype .Type $structs}}{{else}}{{bindtype .Type $structs}}{{end}}
		{{- end}}
	}
	func (e {{errortypename $contract.Type .Normalized.Name}}) Name() string { return "{{.Original.Name}}" }
	func (e {{errortypename $contract.Type .Normalized.Name}}) ErrorID() common.Hash { return common.HexToHash("{{.Original.ID}}") }
	func (e {{errortypename $contract.Type .Normalized.Name}}) Unpack(raw []byte) (*{{errortypename $contract.Type .Normalized.Name}}, error) {
		out := new({{errortypename $contract.Type .Normalized.Name}})
		if err := _{{$contract.Type}}ABI.UnpackIntoInterface(out, e.Name(), raw); err != nil {
			return nil, err
		}
			return out, nil
		}
	{{end}}

	{{ if .Events }}
	// UnpackEvent unpacks event log based on topic0.
		func ({{ decapitalise $contract.Type}} *{{$contract.Type}}) UnpackEvent(log *types.Log) (any, error) {
				if len(log.Topics) == 0 {
					return nil, errNoEventSignature
				}
		for _, item := range {{decapitalise $contract.Type}}EventUnpackers {
			if item.id == log.Topics[0] {
				return item.unpack(log)
			}
		}
			return nil, errEventSignatureMismatch
		}
	{{ end }}

	{{ if .Errors }}
	// UnpackError attempts to decode the provided error data using user-defined
	// error definitions.
		func ({{ decapitalise $contract.Type}} *{{$contract.Type}}) UnpackError(raw []byte) (any, error) {
			if len(raw) < 4 {
				return nil, errErrorSignatureMismatch
			}
			var id [4]byte
			copy(id[:], raw[:4])
			for _, item := range {{decapitalise $contract.Type}}ErrorUnpackers {
				if item.id == id {
					return item.unpack(raw[4:])
				}
			}
			return nil, errErrorSignatureMismatch
		}
	{{ end }}

	type {{decapitalise .Type}}EventUnpacker struct {
		id common.Hash
		unpack func(*types.Log) (any, error)
	}

	type {{decapitalise .Type}}ErrorUnpacker struct {
		id [4]byte
		unpack func([]byte) (any, error)
	}

	{{ if .Events }}
	var {{decapitalise .Type}}EventUnpackers = []{{decapitalise .Type}}EventUnpacker{
		{{- range .Events}}
		{id: common.HexToHash("{{.Original.ID}}"), unpack: func(log *types.Log) (any, error) {
			return ({{eventtypename $contract.Type .Normalized.Name}}{}).Unpack(log)
		}},
		{{- end}}
	}
	{{ end }}

	{{ if .Errors }}
	var {{decapitalise .Type}}ErrorUnpackers = []{{decapitalise .Type}}ErrorUnpacker{
		{{- range .Errors}}
		{id: {{hashidbytes .Original.ID}}, unpack: func(raw []byte) (any, error) {
			return ({{errortypename $contract.Type .Normalized.Name}}{}).Unpack(raw)
		}},
		{{- end}}
	}
	{{ end }}

	func {{decapitalise .Type}}Must[T any](val T, err error) T {
		if err != nil {
			panic(err)
		}
		return val
	}

	func {{decapitalise .Type}}ConvertType[T any](value any) T {
		return *abi.ConvertType(value, new(T)).(*T)
	}

	func {{decapitalise .Type}}CallRaw(ctx context.Context, client bind.ContractCaller, addr common.Address, opts *bind.CallOpts, input []byte) ([]byte, error) {
		var callOpts bind.CallOpts
		if opts != nil {
			callOpts = *opts
		}
		callOpts.Context = ctx
		return bind.NewBoundContract(addr, *_{{.Type}}ABI, client, nil, nil).CallRaw(&callOpts, input)
	}

	func {{decapitalise .Type}}UnpackInput[T any](name string, callData []byte) (*T, error) {
		method, ok := _{{.Type}}ABI.Methods[name]
		if !ok {
			return nil, errors.New("method '" + name + "' not found")
		}
		if len(callData) < 4 || !bytes.Equal(callData[:4], method.ID[:4]) {
			return nil, errMethodSignatureMismatch
		}
		arguments, err := method.Inputs.Unpack(callData[4:])
		if err != nil {
			return nil, err
		}
		input := new(T)
		any(input).(interface{ unpack([]any) }).unpack(arguments)
		return input, nil
	}

	func {{decapitalise .Type}}UnpackEvent[T any](name string, log *types.Log) (*T, error) {
		event, ok := _{{.Type}}ABI.Events[name]
		if !ok {
			return nil, errors.New("event '" + name + "' not found")
		}
		if len(log.Topics) == 0 {
			return nil, errNoEventSignature
		}
		if log.Topics[0] != event.ID {
			return nil, errEventSignatureMismatch
		}
		var values []any
		if len(log.Data) > 0 {
			var err error
			values, err = event.Inputs.Unpack(log.Data)
			if err != nil {
				return nil, err
			}
		}
		out := new(T)
		eventData := any(out).(interface {
			unpack([]any)
			setRaw(*types.Log)
		})
		eventData.unpack(values)
		var indexed abi.Arguments
		for _, arg := range event.Inputs {
			if arg.Indexed {
				indexed = append(indexed, arg)
			}
		}
		if err := abi.ParseTopics(out, indexed, log.Topics[1:]); err != nil {
			return nil, err
		}
		eventData.setRaw(log)
		return out, nil
	}
{{end}}
