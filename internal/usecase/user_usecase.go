package usecase

import (
	"context"
	"fmt"
	"grpc-server-1/internal/entity"
	"grpc-server-1/internal/repository"
	"grpc-server-1/protogen/user"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, log *logrus.Logger, userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            log,
		UserRepository: userRepository,
	}
}

func (u *UserUseCase) GetUser(ctx context.Context, request *user.UserRequest) (*user.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	newUser := new(entity.User)
	if err := u.UserRepository.FindById(tx, newUser, request.Id); err != nil {
		u.Log.WithError(err).Error("error finding user")
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	fmt.Println(newUser)

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, status.Errorf(codes.Internal, "error committing transaction: %v", err)
	}

	return &user.UserResponse{
		Id:         newUser.ID,
		Name:       newUser.Name,
		Occupation: newUser.Occupation,
	}, nil

}
