package usecases

import (
	"context"
	"reliability_demo_app/internal/domain"

	"github.com/melisource/fury_go-core/pkg/log"
)

var _ SearchService = &SearchServiceImpl{}

type SearchService interface {
	SearchItemsByTerm(ctx context.Context, term string) ([]domain.Item, error)
}

func NewSearchService(s ItemSearcher) *SearchServiceImpl {
	return &SearchServiceImpl{
		s,
	}
}

type SearchServiceImpl struct {
	searcher ItemSearcher
}

func (s *SearchServiceImpl) SearchItemsByTerm(ctx context.Context, term string) ([]domain.Item, error) {
	items, err := s.searcher.SearchByTerm(ctx, term)
	if err != nil {
		log.Error(ctx, "search by term failed", log.String("term", term), log.Err(err))
		return nil, err
	}

	return items, nil
}
