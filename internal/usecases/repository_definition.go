package usecases

import (
	"context"
	"reliability_demo_app/internal/domain"
)

type ItemRepository interface {
	UpdatePrice(ctx context.Context, item domain.Item, newPrice float64) error
	Get(ctx context.Context, itemID string) (domain.Item, error)
}

type PriceChangeRepository interface {
	Create(ctx context.Context, input domain.PriceChangeRequest) error
	GetLastByItem(ctx context.Context, item domain.Item) (*domain.PriceChangeRequest, error)
	AllByItem(ctx context.Context, item domain.Item) ([]domain.PriceChangeRequest, error)
}
