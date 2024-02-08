package usecases

import (
	"context"
	"reliability_demo_app/internal/domain"
)

type ItemSearcher interface {
	SearchByTerm(ctx context.Context, term string) ([]domain.Item, error)
}

type BusinessFlowProcesor[T any] interface {
	BeginFlow(ctx context.Context, input T, callbackURL string) error
}
