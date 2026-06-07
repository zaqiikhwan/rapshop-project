package model

import (
	"errors"
	"testing"
)

func TestValidatePaymentProofUpload(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		size    int64
		sniff   []byte
		wantExt string
		wantErr error
	}{
		{
			name:    "png",
			file:    "proof.PNG",
			size:    128,
			sniff:   []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a},
			wantExt: ".png",
		},
		{
			name:    "jpeg",
			file:    "proof.jpeg",
			size:    128,
			sniff:   []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01},
			wantExt: ".jpeg",
		},
		{
			name:    "heic",
			file:    "proof.heic",
			size:    128,
			sniff:   []byte{0x00, 0x00, 0x00, 0x18, 'f', 't', 'y', 'p', 'h', 'e', 'i', 'c'},
			wantExt: ".heic",
		},
		{
			name:    "empty",
			file:    "proof.png",
			size:    0,
			sniff:   nil,
			wantErr: ErrEmptyPaymentProof,
		},
		{
			name:    "too large",
			file:    "proof.png",
			size:    MaxPaymentProofUploadSize + 1,
			sniff:   []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a},
			wantErr: ErrPaymentProofTooLarge,
		},
		{
			name:    "invalid extension",
			file:    "proof.exe",
			size:    128,
			sniff:   []byte("MZ executable"),
			wantErr: ErrInvalidPaymentProofExtension,
		},
		{
			name:    "extension content mismatch",
			file:    "proof.png",
			size:    128,
			sniff:   []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 0x4a, 0x46, 0x49, 0x46, 0x00, 0x01},
			wantErr: ErrInvalidPaymentProofImageBytes,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExt, err := ValidatePaymentProofUpload(tt.file, tt.size, tt.sniff)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error got %v, want %v", err, tt.wantErr)
			}
			if gotExt != tt.wantExt {
				t.Fatalf("extension got %q, want %q", gotExt, tt.wantExt)
			}
		})
	}
}
