package domain_test

import (
	"golang.org/x/crypto/bcrypt"
	"testing"

	apr "github.com/GehirnInc/crypt/apr1_crypt"
	"github.com/stretchr/testify/assert"

	"github.com/default23/avalanche/domain"
)

func TestClient_ComparePassword(t *testing.T) {
	bcryptPwd, err := bcrypt.GenerateFromPassword([]byte("admin"), 10)
	assert.NoError(t, err)

	aprPwd, err := apr.New().Generate([]byte("admin"), nil)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		hash     string
		password string
		want     bool
	}{
		{
			name:     "bcrypt_valid",
			hash:     string(bcryptPwd),
			password: "admin",
			want:     true,
		},
		{
			name:     "bcrypt_invalid",
			hash:     string(bcryptPwd),
			password: "admin123",
			want:     false,
		},
		{
			name:     "apr_valid",
			hash:     aprPwd,
			password: "admin",
			want:     true,
		},
		{
			name:     "apr_invalid",
			hash:     aprPwd,
			password: "admin123",
			want:     false,
		},
		{
			name:     "sha1_valid",
			hash:     "{SHA}0DPiKuNIrrVmD8IUCuw1hQxNqZc=",
			password: "admin",
			want:     true,
		},
		{
			name:     "sha1_invalid",
			hash:     "{SHA}0DPiKuNIrrVmD8IUCuw1hQxNqZc=",
			password: "admin123",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := domain.Client{Password: tt.hash}
			got := c.ComparePassword(tt.password)

			assert.Equal(t, tt.want, got)
		})
	}
}
