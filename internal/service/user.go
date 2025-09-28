package service

import (
	"context"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database/model"
	"gorm.io/gorm"
)

type UserService interface {
	FindByID(ctx context.Context, id uint) (*model.User, error)
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db database.DatabaseService) UserService {
	return &userService{db: db.DB()}
}

func (s *userService) FindByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	err := s.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
