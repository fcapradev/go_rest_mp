package integrations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reliability_demo_app/internal/domain"
	"reliability_demo_app/internal/httpclient"
	"reliability_demo_app/internal/usecases"
	"reliability_demo_app/internal/utils"

	"github.com/melisource/fury_go-core/pkg/rusty"
)

var _ usecases.ItemSearcher = &SearchEngine{}

func NewSearchEngine(factory httpclient.EndpointFactory) *SearchEngine {
	return &SearchEngine{
		getEndpoint: factory.Build("/v1/search"),
	}
}

type SearchEngine struct {
	getEndpoint *rusty.Endpoint
}

func (s *SearchEngine) SearchByTerm(ctx context.Context, term string) ([]domain.Item, error) {
	var query url.Values = make(url.Values)
	query.Add("term", term)
	params := []rusty.RequestOption{
		rusty.WithQuery(query),
	}
	res, err := s.getEndpoint.Get(ctx, params...)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("fetching service from fury api status code not ok: %d", res.StatusCode)
	}

	var response []SearchItemResponse
	if err := json.Unmarshal(res.Body, &response); err != nil {
		return nil, err
	}

	predicate := func(input SearchItemResponse) domain.Item {
		return domain.Item{
			ID:       input.ID,
			Price:    input.Price,
			Name:     input.Name,
			Category: input.Category,
		}
	}

	result, err := utils.Map(response, predicate)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type SearchItemResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}
