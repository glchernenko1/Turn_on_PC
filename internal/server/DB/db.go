package DB

import "Turn_on_PC/internal/DTO"

type DB interface {
	AddUser(user *DTO.User) (uint, error)
	FiendUserByLogin(login string) (*DTO.UserWithID, error)
	DeleteUserByID(id uint) error
}
