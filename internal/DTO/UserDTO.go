package DTO

type BaseUser struct {
	Login string `json:"login" validate:"required,min=4,max=15"`
	Email string `validate:"required,email" json:"email"`
}

type UserSingIn struct {
	Login    string `json:"login" validate:"required,min=4,max=15"`
	Password string `json:"password" validate:"min=8,max=50"`
	Scope    string `json:"scope" validate:"required,in=WebUser,ServerUser,ClientPC"`
}

type CreateUser struct {
	BaseUser
	Password string `json:"password" validate:"min=8,max=50"`
}

type User struct {
	BaseUser
	PasswordHash string `json:"password_hash"`
}

type UserWithID struct {
	ID uint `json:"id"`
	User
}
