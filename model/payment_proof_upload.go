package model

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"
)

const MaxPaymentProofUploadSize = 5 * 1024 * 1024

var (
	ErrEmptyPaymentProof             = errors.New("payment proof file is empty")
	ErrPaymentProofTooLarge          = errors.New("payment proof file exceeds maximum size")
	ErrInvalidPaymentProofExtension  = errors.New("payment proof file extension is not allowed")
	ErrInvalidPaymentProofImageBytes = errors.New("payment proof file content is not a supported image")
)

func ValidatePaymentProofUpload(filename string, size int64, sniff []byte) (string, error) {
	if size <= 0 {
		return "", ErrEmptyPaymentProof
	}
	if size > MaxPaymentProofUploadSize {
		return "", ErrPaymentProofTooLarge
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if !allowedPaymentProofExtensions[ext] {
		return "", ErrInvalidPaymentProofExtension
	}
	if !isAllowedPaymentProofContent(ext, sniff) {
		return "", ErrInvalidPaymentProofImageBytes
	}

	return ext, nil
}

var allowedPaymentProofExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".heic": true,
	".heif": true,
}

func isAllowedPaymentProofContent(ext string, sniff []byte) bool {
	switch ext {
	case ".png":
		return http.DetectContentType(sniff) == "image/png"
	case ".jpg", ".jpeg":
		return http.DetectContentType(sniff) == "image/jpeg"
	case ".heic", ".heif":
		return hasHEIFBrand(sniff)
	default:
		return false
	}
}

func hasHEIFBrand(sniff []byte) bool {
	if len(sniff) < 12 {
		return false
	}
	if string(sniff[4:8]) != "ftyp" {
		return false
	}

	brand := string(sniff[8:12])
	switch brand {
	case "heic", "heix", "hevc", "hevx", "mif1", "msf1":
		return true
	default:
		return false
	}
}
