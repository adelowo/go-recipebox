package hasher

import (
	"golang.org/x/crypto/bcrypt"
)

type BcryptHasher struct {
}

func (b BcryptHasher) Hash(p string) (string, error) {

	hashedInBytes, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashedInBytes), nil

}

func NewBcryptHasher() BcryptHasher {
	return BcryptHasher{}
}
