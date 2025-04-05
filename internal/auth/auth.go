package auth

import "golang.org/x/crypto/bcrypt"

type Provider interface {
	Authenticate() (bool, error)
}

type PasswordProvider struct {
	hash, password string
}

func NewPasswordProvider(hash string, password string) *PasswordProvider {
	return &PasswordProvider{hash: hash, password: password}
}

func (p *PasswordProvider) Authenticate() (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(p.password)); err != nil {
		return false, err
	}
	return true, nil
}

func Authenticate(provider Provider) (bool, error) {
	return provider.Authenticate()
}
