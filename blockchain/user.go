package blockchain

import "crypto/rsa"

type User struct {
	PrivateKey *rsa.PrivateKey
}

func (user *User) Address() string {
	return stringPublic(user.Public())
}

func (user *User) Private() *rsa.PrivateKey {
	return user.PrivateKey
}

func (user *User) Public() *rsa.PublicKey {
	return &user.PrivateKey.PublicKey
}

func (user *User) Purse() string {
	return stringPrivate(user.Private())
}

// New User

const (
	privateKeySize = 512
)

func NewUser() *User {
	return &User{
		PrivateKey: generatePrivate(privateKeySize),
	}
}

func LoadUser(privateData string) *User {
	private := parsePrivate(privateData)
	if private == nil {
		return nil
	}

	return &User{PrivateKey: private}
}
