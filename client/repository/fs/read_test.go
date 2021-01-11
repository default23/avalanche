package repository_test

import (
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/default23/avalanche/client/repository/fs"
	"github.com/default23/avalanche/domain"
)

func Test_Read(t *testing.T) {
	type Suite struct {
		name    string
		login   string
		setup   func(s Suite)
		want    *domain.Client
		wantErr error
	}

	path := ".passwd"
	repo := repository.NewClientRepository(".passwd")

	tests := []Suite{
		{
			name:    "error_open_failed",
			wantErr: &os.PathError{Op: "open", Path: ".passwd", Err: syscall.ENOENT},
		},
		{
			name:  "client_not_found",
			login: "admin",
			setup: func(s Suite) {
				content := `user:user`

				err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
				assert.NoError(t, err)
			},
		},
		{
			name:  "wrong_formatted_credentials",
			login: "admin",
			setup: func(s Suite) {
				content := `
user_user
wrong_format
`

				err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
				assert.NoError(t, err)
			},
		},
		{
			name:  "success",
			login: "admin",
			setup: func(s Suite) {
				content := `
user:user
admin:admin
`
				err := ioutil.WriteFile(path, []byte(content), os.ModePerm)
				assert.NoError(t, err)
			},
			want: &domain.Client{Login: "admin", Password: "admin"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt)
			}

			got, err := repo.Read(tt.login)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)

			_ = os.Remove(path)
		})
	}
}
