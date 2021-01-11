package repository

import (
	"bufio"
	"os"
	"strings"

	"github.com/default23/avalanche/domain"
)

// Read searches the user by login in passwd file
func (r *repo) Read(login string) (*domain.Client, error) {
	passwd, err := os.Open(r.passwdPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = passwd.Close()
	}()

	scanner := bufio.NewScanner(passwd)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		credentials := strings.Split(scanner.Text(), ":")
		if len(credentials) != 2 {
			continue
		}

		if credentials[0] == login {
			return &domain.Client{Login: login, Password: credentials[1]}, nil
		}
	}

	return nil, nil
}
