package DTO

import "Turn_on_PC/pkg/hash"

func Create(user *CreateUser) (*User, error) {
	outUser := new(User)
	outUser.BaseUser = user.BaseUser
	_hash, err := hash.HashPassword(user.Password)
	if err != nil {
		return outUser, err
	}
	outUser.PasswordHash = _hash
	return outUser, err
}
