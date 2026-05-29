package crypto

import (
	"encoding/hex"
	"go/parser"
	"go/token"
	"strconv"
	"testing"
)

func TestKeccak256UsesGoSHA3(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "mapping.go", nil, parser.ImportsOnly)
	if err != nil {
		t.Fatal(err)
	}

	for _, importSpec := range file.Imports {
		importPath, err := strconv.Unquote(importSpec.Path.Value)
		if err != nil {
			t.Fatal(err)
		}
		if importPath == "golang.org/x/crypto/sha3" {
			return
		}
	}

	t.Fatal("Keccak256 should import golang.org/x/crypto/sha3")
}

func TestKeccak256KnownVector(t *testing.T) {
	got := hex.EncodeToString(Keccak256(nil))
	const want = "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	if got != want {
		t.Fatalf("Keccak256(nil) = %s, want %s", got, want)
	}
}
