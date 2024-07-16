package usecase

import (
	"context"
	"fmt"
	"grpc-server-1/internal/entity"
	"grpc-server-1/internal/repository"
	"grpc-server-1/protogen/user"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
		Data: &user.UserData{
			Id:         newUser.ID,
			Name:       newUser.Name,
			Occupation: newUser.Occupation,
			Email:      newUser.Email,
			Role:       newUser.Role,
			CreatedAt:  newUser.CreatedAt.String(), // Assuming CreatedAt is of type time.Time
			UpdatedAt:  newUser.UpdatedAt.String(), // Assuming UpdatedAt is of type time.Time
		},
	}, nil

}

func (u *UserUseCase) CreateUser(ctx context.Context, request *user.UserRequest) (*user.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.PasswordHash), bcrypt.MinCost)
	if err != nil {
		u.Log.WithError(err).Error("error hashing password")
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	newUser := &entity.User{
		ID:           uuid.New().String(),
		Name:         request.Name,
		Occupation:   request.Occupation,
		Email:        request.Email,
		PasswordHash: string(hashedPassword),
		Role:         request.Role,
	}

	totalName, err := u.UserRepository.CountByName(tx, newUser)
	if err != nil {
		u.Log.WithError(err).Error("error checking name availability")
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if totalName > 0 {
		u.Log.Warnf("Name already taken : %+v", err)
		return nil, status.Errorf(codes.Internal, "conflict error")
	}

	totalEmail, err := u.UserRepository.CountByEmail(tx, newUser)
	if err != nil {
		u.Log.WithError(err).Error("error checking email availability")
		return nil, status.Errorf(codes.Internal, "conflict error")
	}

	if totalEmail > 0 {
		u.Log.Warnf("Email already taken : %+v", err)
		return nil, status.Errorf(codes.Internal, "conflict error")
	}

	if err := u.UserRepository.Create(tx, newUser); err != nil {
		u.Log.WithError(err).Error("error creating user")
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	if err := tx.Commit().Error; err != nil {
		u.Log.WithError(err).Error("error committing transaction")
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &user.UserResponse{
		Data: &user.UserData{
			Id:         newUser.ID,
			Name:       newUser.Name,
			Occupation: newUser.Occupation,
			Email:      newUser.Email,
			Role:       newUser.Role,
			CreatedAt:  newUser.CreatedAt.String(), // Assuming CreatedAt is of type time.Time
			UpdatedAt:  newUser.UpdatedAt.String(), // Assuming UpdatedAt is of type time.Time
		},
	}, nil
}
