package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/cache"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database"
	"github.com/sagarmaheshwary/go-microservice-boilerplate/internal/database/model"
	"gorm.io/gorm"
)

type UserService interface {
	FindByID(ctx context.Context, id uint) (*model.User, error)
}

type userService struct {
	database *gorm.DB
	cache    cache.CacheService
}

type UserServiceOpts struct {
	Database database.DatabaseService
	Cache    cache.CacheService
}

func NewUserService(opts *UserServiceOpts) UserService {
	return &userService{database: opts.Database.DB(), cache: opts.Cache}
}

func (s *userService) FindByID(ctx context.Context, id uint) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)

	// Try cache
	if cached, err := s.cache.Get(ctx, cacheKey); err == nil && cached != "" {
		var u model.User
		json.Unmarshal([]byte(cached), &u)
		return &u, nil
	}

	u := &model.User{}
	if err := s.database.First(u, id).Error; err != nil {
		return nil, err
	}

	data, _ := json.Marshal(u)
	s.cache.Set(ctx, cacheKey, data, time.Minute) // cache for 1 minute

	return u, nil
}
