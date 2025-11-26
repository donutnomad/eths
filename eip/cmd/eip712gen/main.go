package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// FieldDefinition 表示结构体字段定义
type FieldDefinition struct {
	Name     string
	Type     string
	GoType   string
	JSONTag  string
	Accessor string
}

// StructDefinition 表示结构体定义
type StructDefinition struct {
	Name   string
	Fields []FieldDefinition
}

// CodeGenConfig 代码生成配置
type CodeGenConfig struct {
	PackageName string
	StructDef   StructDefinition
}

// detectPackageName 自动检测当前目录的包名
func detectPackageName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "main"
	}

	// 查找当前目录下的 .go 文件
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil || len(files) == 0 {
		return "main"
	}

	// 解析第一个非测试文件获取包名
	fset := token.NewFileSet()
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(fset, file, nil, parser.PackageClauseOnly)
		if err != nil {
			continue
		}
		if f.Name != nil && f.Name.Name != "" {
			return f.Name.Name
		}
	}

	return "main"
}

func main() {
	var (
		definition = flag.String("def", "", "合约结构体定义字符串，例如: 'SwapData(address token,address nft,uint256 nftId,...)'")
		output     = flag.String("o", "", "输出文件路径，不指定则输出到标准输出")
		structOnly = flag.Bool("struct-only", false, "仅生成结构体，不生成EIP712函数")
	)
	flag.Parse()

	if *definition == "" {
		fmt.Fprintf(os.Stderr, "错误: 必须提供结构体定义 (-def)\n")
		fmt.Fprintf(os.Stderr, "使用方法: %s -def 'SwapData(address token,address nft,uint256 nftId,...)'\n", os.Args[0])
		flag.Usage()
		os.Exit(1)
	}

	// 自动检测包名
	packageName := detectPackageName()

	// 解析结构体定义
	structDef, err := parseStructDefinition(*definition)
	if err != nil {
		fmt.Fprintf(os.Stderr, "解析错误: %v\n", err)
		os.Exit(1)
	}

	// 生成代码
	config := CodeGenConfig{
		PackageName: packageName,
		StructDef:   *structDef,
	}

	var code string
	if *structOnly {
		code, err = generateStructCode(config)
	} else {
		code, err = generateFullCode(config)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "生成代码错误: %v\n", err)
		os.Exit(1)
	}

	// 输出代码
	if *output != "" {
		err = WriteFormat(*output, []byte(code))
		if err != nil {
			fmt.Fprintf(os.Stderr, "写入文件错误: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("代码已生成到: %s\n", *output)
	} else {
		fmt.Print(code)
	}
}

// parseStructDefinition 解析合约结构体定义字符串
func parseStructDefinition(definition string) (*StructDefinition, error) {
	// 正则表达式匹配结构体定义
	re := regexp.MustCompile(`^(\w+)\(([^)]*)\)$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(definition))
	if len(matches) != 3 {
		return nil, fmt.Errorf("无效的合约结构体定义格式: %s", definition)
	}

	structName := matches[1]
	fieldsStr := matches[2]

	var fields []FieldDefinition
	if strings.TrimSpace(fieldsStr) != "" {
		// 分割字段
		fieldParts := strings.Split(fieldsStr, ",")
		for _, fieldPart := range fieldParts {
			fieldPart = strings.TrimSpace(fieldPart)
			// 分割类型和名称
			parts := strings.Fields(fieldPart)
			if len(parts) != 2 {
				return nil, fmt.Errorf("无效的字段定义: %s", fieldPart)
			}

			solidityType := parts[0]
			fieldName := parts[1]

			// 转换为 Go 类型
			goType, accessor, err := solidityToGoType(solidityType)
			if err != nil {
				return nil, fmt.Errorf("不支持的 Solidity 类型 %s: %v", solidityType, err)
			}

			fields = append(fields, FieldDefinition{
				Name:     fieldName,
				Type:     solidityType,
				GoType:   goType,
				JSONTag:  fieldName,
				Accessor: accessor,
			})
		}
	}

	return &StructDefinition{
		Name:   structName,
		Fields: fields,
	}, nil
}

// solidityToGoType 将 Solidity 类型转换为 Go 类型
func solidityToGoType(solidityType string) (goType, accessor string, err error) {
	switch {
	case solidityType == "address":
		return "common.Address", "%s.Hex()", nil
	case solidityType == "string":
		return "string", "%s", nil
	case solidityType == "bool":
		return "bool", "%s", nil
	case solidityType == "bytes":
		return "[]byte", "hexutil.Encode(%s)", nil
	case strings.HasPrefix(solidityType, "bytes") && len(solidityType) > 5:
		// bytes32, bytes8, etc.
		return "[" + solidityType[5:] + "]byte", "hexutil.Encode(%s[:])", nil
	case solidityType == "uint" || solidityType == "uint256":
		return "*big.Int", "%s.String()", nil
	case strings.HasPrefix(solidityType, "uint"):
		// uint8, uint16, uint32, uint64
		bitSize := solidityType[4:]
		switch bitSize {
		case "8":
			return "uint8", "fmt.Sprintf(\"%%d\", %s)", nil
		case "16":
			return "uint16", "fmt.Sprintf(\"%%d\", %s)", nil
		case "32":
			return "uint32", "fmt.Sprintf(\"%%d\", %s)", nil
		case "64":
			return "uint64", "fmt.Sprintf(\"%%d\", %s)", nil
		default:
			return "*big.Int", "%s.String()", nil
		}
	case solidityType == "int" || solidityType == "int256":
		return "*big.Int", "%s.String()", nil
	case strings.HasPrefix(solidityType, "int"):
		// int8, int16, int32, int64
		bitSize := solidityType[3:]
		switch bitSize {
		case "8":
			return "int8", "fmt.Sprintf(\"%%d\", %s)", nil
		case "16":
			return "int16", "fmt.Sprintf(\"%%d\", %s)", nil
		case "32":
			return "int32", "fmt.Sprintf(\"%%d\", %s)", nil
		case "64":
			return "int64", "fmt.Sprintf(\"%%d\", %s)", nil
		default:
			return "*big.Int", "%s.String()", nil
		}
	case strings.HasSuffix(solidityType, "[]"):
		// 数组类型
		elementType := solidityType[:len(solidityType)-2]
		elementGoType, _, err := solidityToGoType(elementType)
		if err != nil {
			return "", "", err
		}
		return "[]" + elementGoType, "/* TODO: 数组类型访问器 */", nil
	default:
		return "", "", fmt.Errorf("不支持的类型: %s", solidityType)
	}
}

// needsHexutil 检查是否需要导入 hexutil
func needsHexutil(fields []FieldDefinition) bool {
	for _, field := range fields {
		if field.Type == "bytes" || strings.HasPrefix(field.Type, "bytes") && len(field.Type) > 5 {
			return true
		}
	}
	return false
}

// buildImports 构建导入语句
func buildImports(needsHexutil bool, isFullCode bool) string {
	imports := []string{
		`"math/big"`,
	}

	if isFullCode {
		imports = append(imports,
			``,
			`"github.com/donutnomad/eths/eip/eip712"`,
			`"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"`,
		)
	}

	imports = append(imports, `"github.com/ethereum/go-ethereum/common"`)

	if needsHexutil {
		imports = append(imports, `"github.com/ethereum/go-ethereum/common/hexutil"`)
	}

	if isFullCode {
		imports = append(imports, `"github.com/ethereum/go-ethereum/signer/core/apitypes"`)
	}

	return strings.Join(imports, "\n\t")
}

// generateStructCode 仅生成结构体代码
func generateStructCode(config CodeGenConfig) (string, error) {
	tmpl := `package {{.PackageName}}

import (
	` + buildImports(needsHexutil(config.StructDef.Fields), false) + `
)

type {{.StructDef.Name}} struct {
{{- range .StructDef.Fields}}
	{{.Name | title}} {{.GoType}} ` + "`json:\"{{.JSONTag}}\"`" + `
{{- end}}
}
`

	t, err := template.New("struct").Funcs(template.FuncMap{
		"title": strings.Title,
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = t.Execute(&buf, config)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateFullCode 生成完整代码（结构体 + EIP712函数）
func generateFullCode(config CodeGenConfig) (string, error) {
	tmpl := `// Code generated by eip712gen. DO NOT EDIT.

package {{.PackageName}}

import (
	` + buildImports(needsHexutil(config.StructDef.Fields), true) + `
)

type {{.StructDef.Name}} struct {
{{- range .StructDef.Fields}}
	{{.Name | title}} {{.GoType}} ` + "`json:\"{{.JSONTag}}\"`" + `
{{- end}}
}

func (data *{{.StructDef.Name}}) ToMessage() map[string]any {
	return apitypes.TypedDataMessage{
{{- range .StructDef.Fields}}
		"{{.JSONTag}}": {{printf .Accessor (printf "data.%s" (.Name | title))}},
{{- end}}
	}
}

func Generate{{.StructDef.Name}}HashWith(contract common.Address, client bind.ContractCaller, data {{.StructDef.Name}}) ([]byte, apitypes.TypedData, error) {
	domain, err := eip712.GetEIP712Domain(contract, client)
	if err != nil {
		return nil, apitypes.TypedData{}, err
	}
	return Generate{{.StructDef.Name}}Hash(*domain, data)
}

func Generate{{.StructDef.Name}}Hash(domain eip712.Eip712DomainOutput, data {{.StructDef.Name}}) ([]byte, apitypes.TypedData, error) {
	domainTypes, dataDomain := eip712.GetDomainTypes(domain)
	typedData := apitypes.TypedData{
		Types: apitypes.Types{
			"EIP712Domain": domainTypes,
			"{{.StructDef.Name}}": []apitypes.Type{
{{- range .StructDef.Fields}}
				{Name: "{{.JSONTag}}", Type: "{{.Type}}"},
{{- end}}
			},
		},
		PrimaryType: "{{.StructDef.Name}}",
		Domain:      dataDomain,
		Message: data.ToMessage(),
	}
	hash, _, err := apitypes.TypedDataAndHash(typedData)
	if err != nil {
		return nil, apitypes.TypedData{}, err
	}
	return hash, typedData, nil
}
`

	t, err := template.New("full").Funcs(template.FuncMap{
		"title": strings.Title,
	}).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	err = t.Execute(&buf, config)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
