package service

import (
	"context"

	"github.com/VikaPaz/pantheon/internal/models"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	repo Repository
	log  *logrus.Logger
}

type Repository interface {
	Create(user models.User) (models.User, error)
	GetById(user models.User) (models.User, error)
	GetByUsername(user models.User) (models.User, error)
	Delete(request models.User) error
}

func NewService(repo Repository, logger *logrus.Logger) *UserService {
	return &UserService{
		repo: repo,
		log:  logger,
	}
}

func (u *UserService) CreateUser(ctx *context.Context, user models.User) (models.User, error) {
	u.log.Debugf("Checking user exists")
	result, err := u.repo.GetByUsername(user)
	if err != nil {
		return models.User{}, err
	}
	if result.Id != "" {
		return models.User{}, models.ErrUserExists
	}

	u.log.Debugf("Creating user: %v", user)
	userInf, err := u.repo.Create(user)
	if err != nil {
		return models.User{}, err
	}

	return userInf, nil
}

func (u *UserService) DeleteUser(ctx *context.Context, user models.User) error {
	u.log.Debugf("Deleting user with ID: %v", user.Id)
	err := u.repo.Delete(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetByUsername(ctx *context.Context, user models.User) (models.User, error) {
	u.log.Debugf("Getting users with filter: %v", user)
	result, err := u.repo.GetByUsername(user)
	if err != nil {
		return models.User{}, err
	}
	return result, nil
}

func (u *UserService) GetById(ctx *context.Context, user models.User) (models.User, error) {
	u.log.Debugf("Getting users with filter: %v", user)
	result, err := u.repo.GetById(user)
	if err != nil {
		return models.User{}, err
	}
	return result, nil
}
