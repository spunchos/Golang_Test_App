package user

import (
	"TestApp/pkg/logging"
	"context"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) *Service {
	return &Service{storage: storage, logger: logger}
}

func (s *Service) Create(ctx context.Context, user *User) (User, error) {
	s.logger.Trace(user)
	id, err := s.storage.Create(ctx, user)
	if err != nil {
		s.logger.Errorf("failed to create user %v", err)
		return User{}, err
	}
	s.logger.Tracef("created user with id %s", id)
	user.ID = id
	return *user, nil
}

func (s *Service) FindOne(ctx context.Context, id string) (User, error) {
	user, err := s.storage.FindOne(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to find user %v", err)
		return User{}, err
	}
	s.logger.Tracef("found user with id %s", id)
	return user, nil
}

func (s *Service) FindAll(ctx context.Context) ([]User, error) {
	users, err := s.storage.FindAll(ctx)
	if err != nil {
		s.logger.Errorf("failed to find users %v", err)
		return nil, err
	}
	s.logger.Tracef("found %d users", len(users))
	return users, nil
}

func (s *Service) Update(ctx context.Context, user User) error {
	err := s.storage.Update(ctx, user)
	if err != nil {
		s.logger.Errorf("failed to update user %v", err)
		return err
	}
	s.logger.Tracef("updated user with id %s", user.ID)
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		s.logger.Errorf("failed to delete user %v", err)
		return err
	}
	s.logger.Tracef("deleted user with id %s", id)
	return nil
}
