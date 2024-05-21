package server

import "context"

type Repository interface {
	Add(ctx context.Context, log Log) error
	GetAll(ctx context.Context) ([]Log, error)
}

type Service interface {
	Add(ctx context.Context, log Log) error
	GetAll(ctx context.Context) ([]Log, error)
}

type ServiceMiddleware func(Service) Service
