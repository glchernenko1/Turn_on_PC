package postgres

import (
	"Turn_on_PC/pkg/logging"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"Turn_on_PC/internal/DTO"
	"encoding/json"
	"Turn_on_PC/internal/server/config"
	"Turn_on_PC/internal/server/DB"
)

type db struct {
	Logger *logging.Logger
	DB     *gorm.DB
}

func NewDB(cfg *config.Config, logger *logging.Logger) DB.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		cfg.Postgres.Host, cfg.Postgres.Username, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Port)
	logger.Info(dsn)
	_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("DB connected to database")
	err = _db.AutoMigrate(&User{})
	if err != nil {
		logger.Fatal("AutoMigrate failed")
	}
	logger.Info("DB successfully migrated")

	return &db{Logger: logger, DB: _db}
}

func (db *db) AddUser(user *DTO.User) (uint, error) {
	userJson, err := json.Marshal(user)
	if err != nil {
		return 0, err
	}
	userDB := new(User)
	err = json.Unmarshal([]byte(userJson), &userDB)
	if err != nil {
		return 0, err
	}

	tmp := db.DB.Create(userDB)
	err = tmp.Error
	return userDB.ID, err
}

func (db *db) FiendUserByLogin(login string) (*DTO.UserWithID, error) {
	var user User
	result := db.DB.Where("login = ?", login).Find(&user)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("not found")
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	userOutput := new(DTO.UserWithID)
	err = json.Unmarshal([]byte(userJson), &userOutput)
	return userOutput, err
}

func (db *db) DeleteUserByID(id uint) error {
	return db.DB.Delete(&User{}, id).Error
}
