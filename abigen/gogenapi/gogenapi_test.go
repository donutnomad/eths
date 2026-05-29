package gogenapi

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadInputArtifactAndGenerate(t *testing.T) {
	dir := t.TempDir()
	artifact := filepath.Join(dir, "escrow_pack.json")
	err := os.WriteFile(artifact, []byte(`{
		"_format": "hh-sol-artifact-1",
		"contractName": "Escrow",
		"abi": [{"type":"function","name":"coordinator","inputs":[],"outputs":[{"type":"address"}],"stateMutability":"view"}],
		"bytecode": "0x60016002"
	}`), 0600)
	if err != nil {
		t.Fatal(err)
	}

	input, err := LoadInput(artifact)
	if err != nil {
		t.Fatal(err)
	}
	if input.Kind != InputKindArtifact {
		t.Fatalf("kind = %v, want artifact", input.Kind)
	}
	if input.ContractName != "Escrow" {
		t.Fatalf("contractName = %q, want Escrow", input.ContractName)
	}
	if DefaultTypeName(input) != "Escrow" {
		t.Fatalf("default type = %q, want Escrow", DefaultTypeName(input))
	}

	code, err := Generate(GenerateOptions{
		ABIPath:     artifact,
		PackageName: "aaa",
	})
	if err != nil {
		t.Fatal(err)
	}
	output := string(code)
	for _, want := range []string{
		"package aaa",
		"type Escrow struct",
		`Bin: "0x60016002"`,
		`github.com/ethereum/go-ethereum/accounts/abi/bind/v2`,
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("generated code missing %q", want)
		}
	}
}

func TestDefaultTypeNamePlainABIUsesFileName(t *testing.T) {
	dir := t.TempDir()
	abiFile := filepath.Join(dir, "escrow_pack.json")
	err := os.WriteFile(abiFile, []byte(`[{"type":"function","name":"approve","inputs":[],"outputs":[]}]`), 0600)
	if err != nil {
		t.Fatal(err)
	}

	input, err := LoadInput(abiFile)
	if err != nil {
		t.Fatal(err)
	}
	if input.Kind != InputKindABI {
		t.Fatalf("kind = %v, want abi", input.Kind)
	}
	if DefaultTypeName(input) != "EscrowPack" {
		t.Fatalf("default type = %q, want EscrowPack", DefaultTypeName(input))
	}
}
