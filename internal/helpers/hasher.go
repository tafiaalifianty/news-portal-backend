package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashAndSalt(pwd string) (string, error)
	ComparePasswords(hashedPwd string, plainPwd []byte) bool
}

type defaultHasher struct{}

func NewHasher() Hasher {
	return &defaultHasher{}
}

func (h *defaultHasher) HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (h *defaultHasher) ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), plainPwd)

	return err == nil
}
