package domain

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	apr "github.com/GehirnInc/crypt/apr1_crypt"
	"golang.org/x/crypto/bcrypt"
)

// hash length allows to identify the
// used algorithm in provided hash string
const (
	bcryptLength = 60
	aprLength    = 37
	shaLength    = 33
)

// ClientRepository is a repo, that stores
// the client credentials for proxy
type ClientRepository interface {
	Read(login string) (*Client, error)
}

// Client is a struct that stores the
// client login and hashed password
type Client struct {
	Login    string
	Password string
}

func (c Client) ComparePassword(pwd string) bool {
	password := []byte(pwd)
	switch len(c.Password) {
	case bcryptLength:
		err := bcrypt.CompareHashAndPassword([]byte(c.Password), password)
		if err != nil {
			return false
		}
		return true

	case aprLength:
		a := apr.New()
		return a.Verify(c.Password, password) == nil

	case shaLength:
		h := sha1.New()
		h.Write(password)
		hashed := base64.StdEncoding.EncodeToString(h.Sum(nil))
		return fmt.Sprintf(`{SHA}%s`, hashed) == c.Password
	}

	return false
}
