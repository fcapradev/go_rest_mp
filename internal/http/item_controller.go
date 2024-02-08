package http

import (
	"net/http"
	"reliability_demo_app/internal/domain"
	"reliability_demo_app/internal/usecases"
	"reliability_demo_app/internal/utils"

	"github.com/google/uuid"
	"github.com/melisource/fury_go-core/pkg/log"
	"github.com/melisource/fury_go-core/pkg/web"
)

var _ Controller = &ItemController{}

func NewItemController(u usecases.ItemService) *ItemController {
	return &ItemController{
		u,
	}
}

type ItemController struct {
	usecase usecases.ItemService
}

func (c *ItemController) AddRoutes(router Router) {
	router.Get("/items", c.getItems)
	router.Get("/items/{id}", c.getItem)
	router.Post("/items/{id}/price-changes", c.createPriceChange)
	router.Get("/items/{id}/price-changes", c.getPriceChanges)
	router.Post("/items/notifications", c.notifications)
}

func (c *ItemController) getItems(w http.ResponseWriter, r *http.Request) error {

	return web.EncodeJSON(w, nil, http.StatusNotImplemented)
}

func (c *ItemController) getItem(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := web.Params(r)
	itemID, err := params.String("id")
	if err != nil {
		log.Error(ctx, "missing id param in the uri", log.Err(err))
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	item, err := c.usecase.GetItem(ctx, itemID)
	if err != nil {
		log.Error(ctx, "get item failed", log.Err(err))
		return web.NewError(http.StatusInternalServerError, err.Error())
	}

	result := GetItemResponse{
		item.ID,
		item.Name,
		item.Price,
	}

	return web.EncodeJSON(w, result, http.StatusOK)
}

type GetItemResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func (c *ItemController) createPriceChange(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := web.Params(r)
	itemID, err := params.String("id")
	if err != nil {
		log.Error(ctx, "missing id param in the uri", log.Err(err))
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	var payload PriceChangeRequest
	if err := web.DecodeJSON(r, &payload); err != nil {
		log.Error(ctx, "unable to unmarshal payload", log.Err(err))
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	result := domain.PriceChangeRequest{
		ItemID:        itemID,
		TransactionID: uuid.NewString(),
		Price:         payload.Price,
	}

	err = c.usecase.SendPriceChangeRequest(ctx, result)
	if err != nil {
		log.Error(ctx, "send price chage request failed", log.Err(err))
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	return web.EncodeJSON(w, nil, http.StatusNoContent)
}

func (c *ItemController) getPriceChanges(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	params := web.Params(r)
	itemID, err := params.String("id")
	if err != nil {
		return web.NewError(http.StatusBadRequest, err.Error())
	}

	item := domain.Item{
		ID: itemID,
	}

	result, err := c.usecase.GetPriceChanges(ctx, item)
	if err != nil {
		return web.NewError(http.StatusInternalServerError, err.Error())
	}

	response, err := utils.Map(result, func(input domain.PriceChangeRequest) PriceChangeResponse {
		return PriceChangeResponse{
			TransactionID: input.TransactionID,
			ItemID:        input.ItemID,
			Price:         input.Price,
		}
	})
	if err != nil {
		return web.NewError(http.StatusInternalServerError, err.Error())
	}

	return web.EncodeJSON(w, response, http.StatusOK)
}

type PriceChangeResponse struct {
	TransactionID string  `json:"transaction_id"`
	ItemID        string  `json:"item_id"`
	Price         float64 `json:"price"`
}

type PriceChangeRequest struct {
	TransactionID string  `json:"transaction_id"`
	Price         float64 `json:"price"`
}

func (c *ItemController) notifications(w http.ResponseWriter, r *http.Request) error {
	var input PriceChangeNotificationMessage
	if err := web.DecodeJSON(r, &input); err != nil {
		return web.EncodeJSON(w, err, http.StatusBadRequest)
	}

	if err := c.usecase.UpdatePrice(r.Context(), input.ItemID); err != nil {
		log.Error(r.Context(), "error creating price change request", log.Err(err))
		return web.EncodeJSON(w, "error creating price change request", http.StatusInternalServerError)
	}

	return web.EncodeJSON(w, nil, http.StatusAccepted)
}

type PriceChangeNotificationMessage struct {
	TransactionID string `json:"transaction_id"`
	ItemID        string `json:"item_id"`
	Status        string `json:"status"`
}
