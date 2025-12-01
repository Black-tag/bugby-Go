package utils

import (
	"log/slog"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	logger := slog.Default().With(
		"test", "HashPasswordTests",
	)
	logger.Info("started tests")

	password := "supersecret"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error but got %v", err)
	}

	if hash == "" {
		t.Fatal("hash is not supposed to be empty")
	}
	if hash == password {
		t.Fatal("hash should not be equal to password")
	}
	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expectrd no error got %v", err)
	}

	if hash == hash2 {
		t.Fatal("hashing for 2 nd time shouldnt give same hashed password")
	}
}

func TestHashPassword2(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "mypassword", false},
		{"empty password", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr = %v", err, tt.wantErr)
				return

			}
			if !tt.wantErr {
				err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(tt.password))
				if err != nil {
					t.Errorf("hashed password does not match the original %v", err)
				}
			}
		})
	}
}
