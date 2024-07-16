package repository

import (
	"grpc-server-1/internal/entity"
	"grpc-server-1/protogen/user"

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

func (r *UserRepository) Search(db *gorm.DB, request *user.UserSearchRequest) ([]entity.User, int64, error) {
	var users []entity.User
	page := int(request.Page)
	size := int(request.Size)

	if err := db.Scopes(r.FilterUser(request)).Offset((page - 1) * size).Limit(size).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.User{}).Scopes(r.FilterUser(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) FilterUser(request *user.UserSearchRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if name := request.Name; name != "" {
			name = "%" + name + "%"
			tx = tx.Where("name LIKE ?", name)
		}

		if email := request.Email; email != "" {
			email = "%" + email + "%"
			tx = tx.Where("email LIKE ?", email)
		}

		if role := request.Role; role != "" {
			tx = tx.Where("role = ?", role)
		}

		return tx
	}
}
