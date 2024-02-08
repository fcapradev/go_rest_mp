package integrations

import (
	"context"
	"encoding/json"
	"net/http"
	"reliability_demo_app/internal/domain"
	"reliability_demo_app/internal/httpclient"
	"reliability_demo_app/internal/usecases"
	"reliability_demo_app/internal/utils"
	"time"

	"github.com/melisource/fury_go-core/pkg/rusty"
	"github.com/melisource/fury_go-core/pkg/web"
)

var (
	_ usecases.ItemRepository        = &ItemsApi{}
	_ usecases.PriceChangeRepository = &ItemsApi{}
)

func NewItemsApi(factory httpclient.EndpointFactory) *ItemsApi {
	return &ItemsApi{
		itemEndpoint:             factory.Build("/v1/products/{id}"),
		priceChangesEndpoint:     factory.Build("/v1/products/{id}/price-changes"),
		priceChangesLastEndpoint: factory.Build("/v1/products/{id}/price-changes/last"),
		postEndpoint:             factory.Build("/v1/price-changes"),
	}
}

type ItemsApi struct {
	postEndpoint             *rusty.Endpoint
	itemEndpoint             *rusty.Endpoint
	priceChangesEndpoint     *rusty.Endpoint
	priceChangesLastEndpoint *rusty.Endpoint
}

func (i *ItemsApi) GetLastByItem(ctx context.Context, item domain.Item) (*domain.PriceChangeRequest, error) {
	res, err := i.priceChangesLastEndpoint.Get(ctx, rusty.WithParam("id", item.ID))
	if err != nil {
		err = web.NewError(res.StatusCode, err.Error())
		return nil, err
	}
	if res.StatusCode != 200 {
		err = web.NewError(res.StatusCode, err.Error())
		return nil, err
	}

	var response LastItemPriceApprovalResponse
	if err := json.Unmarshal(res.Body, &response); err != nil {
		err = web.NewError(http.StatusInternalServerError, err.Error())
		return nil, err
	}

	output := domain.PriceChangeRequest{
		TransactionID: response.TransactionId,
		ItemID:        response.ProductId,
		Price:         response.Price,
	}

	return &output, nil
}

func (i *ItemsApi) AllByItem(ctx context.Context, item domain.Item) ([]domain.PriceChangeRequest, error) {
	res, err := i.priceChangesEndpoint.Get(ctx, rusty.WithParam("id", item.ID))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, err
	}

	var response []ChangePriceResponse
	if err := json.Unmarshal(res.Body, &response); err != nil {
		return nil, err
	}

	output, err := utils.Map(response, func(input ChangePriceResponse) domain.PriceChangeRequest {
		return domain.PriceChangeRequest{
			TransactionID: input.TransactionID,
			ItemID:        input.ItemID,
			Price:         input.Price,
		}
	})
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (a *ItemsApi) Create(ctx context.Context, input domain.PriceChangeRequest) error {
	body := PriceChangeRequestBody{
		TransactionID: input.TransactionID,
		ItemID:        input.ItemID,
		Price:         input.Price,
	}
	params := []rusty.RequestOption{
		rusty.WithBody(body),
		rusty.WithParam("id", input.ItemID),
	}

	res, err := a.priceChangesEndpoint.Post(ctx, params...)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return err
	}

	return nil
}

type PriceChangeRequestBody struct {
	TransactionID string  `json:"transaction_id"`
	ItemID        string  `json:"item_id"`
	Price         float64 `json:"price"`
}

type ChangePriceResponse struct {
	TransactionID string    `json:"transaction_id"`
	ItemID        string    `json:"item_id"`
	Price         float64   `json:"price"`
	CreatedAt     time.Time `json:"created_at"`
}

func (i *ItemsApi) UpdatePrice(ctx context.Context, item domain.Item, newPrice float64) error {
	body := ItemNewPriceRequest{
		Price: newPrice,
	}

	res, err := i.itemEndpoint.Patch(ctx, rusty.WithParam("id", item.ID), rusty.WithBody(body))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return err
	}
	return nil
}

type LastItemPriceApprovalResponse struct {
	ProductId     string  `json:"product_id"`
	TransactionId string  `json:"transaction_id"`
	Price         float64 `json:"price"`
}

type ItemNewPriceRequest struct {
	Price float64 `json:"price"`
}

func (i *ItemsApi) Get(ctx context.Context, itemID string) (domain.Item, error) {
	res, err := i.itemEndpoint.Get(ctx, rusty.WithParam("id", itemID))
	if err != nil {
		return domain.Item{}, err
	}

	if res.StatusCode != 200 {
		return domain.Item{}, err
	}

	var response SearchItemResponse
	if err := json.Unmarshal(res.Body, &response); err != nil {
		return domain.Item{}, err
	}

	output := domain.Item{
		ID:    response.ID,
		Name:  response.Name,
		Price: response.Price,
	}

	return output, nil
}
