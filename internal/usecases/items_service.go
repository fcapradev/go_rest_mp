package usecases

import (
	"context"
	"fmt"
	"reliability_demo_app/internal/domain"

	"github.com/melisource/fury_go-core/pkg/log"
)

var _ ItemService = &ItemServiceImpl{}

type ItemService interface {
	UpdatePrice(ctx context.Context, itemID string) error
	SendPriceChangeRequest(ctx context.Context, input domain.PriceChangeRequest) error
	GetPriceChanges(ctx context.Context, item domain.Item) ([]domain.PriceChangeRequest, error)
	GetItem(ctx context.Context, itemID string) (domain.Item, error)
}

func NewItemService(
	itemRepository ItemRepository,
	priceChangeRepository PriceChangeRepository,
	flowProcessor BusinessFlowProcesor[domain.PriceChangeRequest],
	callbackBaseURL string,
) *ItemServiceImpl {
	return &ItemServiceImpl{
		itemRepository:        itemRepository,
		priceChangeRepository: priceChangeRepository,
		flowProcessor:         flowProcessor,
		callbackURL:           fmt.Sprintf("%s%s", callbackBaseURL, "/items/notifications"),
	}
}

type ItemServiceImpl struct {
	itemRepository        ItemRepository
	priceChangeRepository PriceChangeRepository
	flowProcessor         BusinessFlowProcesor[domain.PriceChangeRequest]
	callbackURL           string
}

func (s *ItemServiceImpl) SendPriceChangeRequest(ctx context.Context, input domain.PriceChangeRequest) error {
	log.Info(ctx, "sending price change request", log.Reflect("input", input))

	go func() {
		if err := s.priceChangeRepository.Create(context.Background(), input); err != nil {
			log.Error(ctx, "create price change failed", log.Err(err))
		}
	}()

	go func() {
		if err := s.flowProcessor.BeginFlow(context.Background(), input, s.callbackURL); err != nil {
			log.Error(ctx, "begin flow failed", log.Err(err))
		}
	}()

	return nil
}

func (s *ItemServiceImpl) UpdatePrice(ctx context.Context, itemID string) error {
	item := domain.Item{ID: itemID}
	priceChangeRequest, err := s.priceChangeRepository.GetLastByItem(ctx, item)
	if err != nil {
		return err
	}

	fmt.Printf("*** priceChangeRequest: %+v\n", priceChangeRequest)

	err = s.itemRepository.UpdatePrice(ctx, item, priceChangeRequest.Price)
	if err != nil {
		return err
	}
	return nil
}

func (s *ItemServiceImpl) GetPriceChanges(ctx context.Context, item domain.Item) ([]domain.PriceChangeRequest, error) {
	result, err := s.priceChangeRepository.AllByItem(ctx, item)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *ItemServiceImpl) GetItem(ctx context.Context, itemID string) (domain.Item, error) {
	result, err := s.itemRepository.Get(ctx, itemID)
	if err != nil {
		return domain.Item{}, err
	}

	return result, nil
}
