package http

import (
	"net/http"
	"reliability_demo_app/internal/domain"
	"reliability_demo_app/internal/usecases"
	"reliability_demo_app/internal/utils"

	"github.com/melisource/fury_go-core/pkg/web"
)

func NewSearchController(u usecases.SearchService) *SearchController {
	return &SearchController{
		u,
	}
}

type SearchController struct {
	usecase usecases.SearchService
}

func (c *SearchController) AddRoutes(router Router) {
	router.Get("/search", c.search)
}

func (c *SearchController) search(w http.ResponseWriter, r *http.Request) error {
	term := r.URL.Query().Get("term")
	items, err := c.usecase.SearchItemsByTerm(r.Context(), term)
	if err != nil {
		return web.EncodeJSON(w, err, http.StatusInternalServerError)
	}

	result, err := utils.Map(items, func(i domain.Item) SearchResponse {
		return SearchResponse{
			ID:    i.ID,
			Name:  i.Name,
			Price: i.Price,
		}
	})
	if err != nil {
		return web.EncodeJSON(w, err, http.StatusInternalServerError)
	}

	return web.EncodeJSON(w, DataResponse{Data: result}, http.StatusOK)
}

type DataResponse struct {
	Data any `json:"data"`
}

type SearchResponse struct {
	Name  string  `json:"name"`
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}
