package contractcall

import (
	"bytes"
	"testing"
)

func TestDecodeHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{
			name:  "empty string",
			input: "",
			want:  []byte{},
		},
		{
			name:  "0x prefix only",
			input: "0x",
			want:  []byte{},
		},
		{
			name:  "0X prefix only",
			input: "0X",
			want:  []byte{},
		},
		{
			name:  "with 0x prefix",
			input: "0x48656c6c6f",
			want:  []byte("Hello"),
		},
		{
			name:  "without prefix",
			input: "48656c6c6f",
			want:  []byte("Hello"),
		},
		{
			name:  "odd length without prefix",
			input: "1",
			want:  []byte{0x01},
		},
		{
			name:  "odd length with prefix",
			input: "0x1",
			want:  []byte{0x01},
		},
		{
			name:    "invalid hex chars",
			input:   "0xZZ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeHex(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("decodeHex(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && !bytes.Equal(got, tt.want) {
				t.Errorf("decodeHex(%q) = %x, want %x", tt.input, got, tt.want)
			}
		})
	}
}
