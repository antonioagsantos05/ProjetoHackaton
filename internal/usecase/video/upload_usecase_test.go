package video

import (
	"testing"
)

// Testa que NewUploadUseCase retorna a interface UploadUseCase corretamente.
func TestNewUploadUseCase_ReturnsInterface(t *testing.T) {
	// NewUploadUseCase aceita db, channel e minioClient.
	// Passamos nil pois estamos apenas validando o contrato da interface.
	uc := NewUploadUseCase(nil, nil, nil)
	if uc == nil {
		t.Fatal("expected UploadUseCase, got nil")
	}
}

// Testa que NewDownloadUseCase retorna a interface DownloadUseCase corretamente.
func TestNewDownloadUseCase_ReturnsInterface(t *testing.T) {
	uc := NewDownloadUseCase(nil, nil)
	if uc == nil {
		t.Fatal("expected DownloadUseCase, got nil")
	}
}

// Testa que NewListUseCase retorna a interface ListUseCase corretamente.
func TestNewListUseCase_ReturnsInterface(t *testing.T) {
	uc := NewListUseCase(nil)
	if uc == nil {
		t.Fatal("expected ListUseCase, got nil")
	}
}
