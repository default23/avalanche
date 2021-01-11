package repository

import "github.com/default23/avalanche/domain"

type repo struct {
	passwdPath string
}

func NewClientRepository(path string) domain.ClientRepository {
	r := new(repo)
	r.passwdPath = path

	return r
}
