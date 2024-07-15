package repository

import (
	"grpc-server-1/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
}

func NewUserRepository(log *logrus.Logger) *UserRepository {
	return &UserRepository{
		Log: log,
	}
}

func (r *UserRepository) CountByName(db *gorm.DB, user *entity.User) (int64, error) {
	var total int64
	err := db.Model(user).Where("name = ?", user.Name).Count(&total).Error
	return total, err
}

func (r *UserRepository) CountByEmail(db *gorm.DB, user *entity.User) (int64, error) {
	var total int64
	err := db.Model(user).Where("email = ?", user.Email).Count(&total).Error
	return total, err
}

func (r *UserRepository) FindByName(db *gorm.DB, user *entity.User, name string) error {
	return db.Where("name = ?", name).Take(user).Error
}

func (r *UserRepository) FindByEmail(db *gorm.DB, user *entity.User, email string) error {
	return db.Where("email = ?", email).Take(user).Error
}
