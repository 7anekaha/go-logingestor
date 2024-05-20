package server

import "context"

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Add(ctx context.Context, log Log) error {
	return s.repo.Add(ctx, log)
}

func (s *service) GetAll(ctx context.Context) ([]Log, error) {
	return s.repo.GetAll(ctx)
}
