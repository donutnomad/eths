package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadABIInputArtifact(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	artifact := filepath.Join(dir, "Escrow2.json")
	err := os.WriteFile(artifact, []byte(`{
		"_format": "hh-sol-artifact-1",
		"contractName": "Escrow",
		"abi": [{"type":"function","name":"coordinator","inputs":[],"outputs":[{"type":"address"}],"stateMutability":"view"}],
		"bytecode": "0x60016002"
	}`), 0600)
	if err != nil {
		t.Fatal(err)
	}

	input, err := readABIInput("@"+artifact, strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(input.ABI), `[{"type":"function","name":"coordinator","inputs":[],"outputs":[{"type":"address"}],"stateMutability":"view"}]`; got != want {
		t.Fatalf("ABI mismatch\ngot:  %s\nwant: %s", got, want)
	}
	if got, want := string(input.Bin), "0x60016002"; got != want {
		t.Fatalf("bin mismatch: got %q, want %q", got, want)
	}
}

func TestReadABIInputPlainFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	abiFile := filepath.Join(dir, "Escrow.abi")
	rawABI := `[{"type":"function","name":"approve","inputs":[],"outputs":[]}]`
	err := os.WriteFile(abiFile, []byte(rawABI), 0600)
	if err != nil {
		t.Fatal(err)
	}

	input, err := readABIInput(abiFile, strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	if got := string(input.ABI); got != rawABI {
		t.Fatalf("ABI mismatch: got %q, want %q", got, rawABI)
	}
	if len(input.Bin) != 0 {
		t.Fatalf("bin mismatch: got %q, want empty", input.Bin)
	}
}

func TestGenerateUsesArtifactBytecode(t *testing.T) {
	dir := t.TempDir()
	artifact := filepath.Join(dir, "Escrow2.json")
	outFile := filepath.Join(dir, "escrow.go")
	err := os.WriteFile(artifact, []byte(`{
		"_format": "hh-sol-artifact-1",
		"contractName": "Escrow",
		"abi": [{"type":"function","name":"coordinator","inputs":[],"outputs":[{"type":"address"}],"stateMutability":"view"}],
		"bytecode": "0x60016002"
	}`), 0600)
	if err != nil {
		t.Fatal(err)
	}

	err = app.Run([]string{
		"abigen",
		"--pkg", "aaa",
		"--type", "EscrowPkg",
		"--abi", "@" + artifact,
		"--out", outFile,
		"--v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	generated, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(generated), `Bin: "0x60016002"`) {
		t.Fatalf("generated binding does not contain artifact bytecode")
	}
}

func TestGenerateBinFlagOverridesArtifactBytecode(t *testing.T) {
	dir := t.TempDir()
	artifact := filepath.Join(dir, "Escrow2.json")
	binFile := filepath.Join(dir, "Escrow.bin")
	outFile := filepath.Join(dir, "escrow.go")
	err := os.WriteFile(artifact, []byte(`{
		"_format": "hh-sol-artifact-1",
		"contractName": "Escrow",
		"abi": [{"type":"function","name":"coordinator","inputs":[],"outputs":[{"type":"address"}],"stateMutability":"view"}],
		"bytecode": "0x60016002"
	}`), 0600)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binFile, []byte("0x60036004"), 0600); err != nil {
		t.Fatal(err)
	}

	err = app.Run([]string{
		"abigen",
		"--pkg", "aaa",
		"--type", "EscrowPkg",
		"--abi", "@" + artifact,
		"--bin", binFile,
		"--out", outFile,
		"--v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	generated, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(generated), `Bin: "0x60036004"`) {
		t.Fatalf("generated binding does not contain --bin bytecode")
	}
	if strings.Contains(string(generated), `Bin: "0x60016002"`) {
		t.Fatalf("generated binding contains artifact bytecode")
	}
}
