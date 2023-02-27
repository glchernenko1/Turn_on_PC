package postgres

type User struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Login        string `gorm:"not null; type:varchar(50); unique;" json:"login"`
	Email        string `gorm:"not null; type:varchar(50); unique;" json:"email"`
	PasswordHash string `gorm:"not null" json:"password_hash"`
}
