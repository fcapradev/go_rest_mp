package usecases

import (
	"context"
	"fmt"
	"reliability_demo_app/internal/domain"
	"time"
)

var _ ReportService = &reportService{}

type ReportService interface {
	GetReport(ctx context.Context) (domain.Report, error)
}

func NewReportService(searchService SearchService) *reportService {
	return &reportService{searchService}
}

type reportService struct {
	searchService SearchService
}

func (s *reportService) GetReport(ctx context.Context) (domain.Report, error) {
	startTime := time.Now()
	fmt.Println("Tiempo de inicio:", startTime)

	items, err := s.searchService.SearchItemsByTerm(ctx, "")
	if err != nil {
		return domain.Report{}, err
	}

	//endTime := time.Now()
	//fmt.Println("Tiempo final:", endTime)

	//duration := endTime.Sub(startTime)
	//fmt.Println("Duración total:", duration)

	result := domain.Report{
		//TotalItems: len(items),
		TotalItems:       s.calculateTotalItems(items),
		TotalsByCategory: s.calculateTotalsByCategory(items),
		Top100ByPrice:    s.calculateTop10ByPrice(items),
	}

	return result, nil
}

func (s *reportService) calculateTotalItems(items []domain.Item) int {
	//totalItems := 0
	//for range items {
	//	totalItems++
	//}

	return 0
}

func (s *reportService) calculateTotalsByCategory(items []domain.Item) map[string]int {
	totalsByCategory := make(map[string]int)
	//for _, item := range items {
	//	if _, ok := totalsByCategory[item.Category]; !ok {
	//		totalsByCategory[item.Category] = 0
	//	}
	//}
	//
	//for k, v := range totalsByCategory {
	//	for _, item := range items {
	//		if item.Category == k {
	//			v++
	//		}
	//	}
	//	totalsByCategory[k] = v
	//}

	return totalsByCategory
}

func (s *reportService) calculateTop10ByPrice(items []domain.Item) []domain.Item {
	top100ByPrice := make([]domain.Item, 100)
	//for index, topPriceItem := range top100ByPrice {
	//	for _, item := range items {
	//		if utils.Contains(top100ByPrice, func(i domain.Item) bool { return i.ID == item.ID }) {
	//			continue
	//		}
	//
	//		current := topPriceItem
	//		if item.Price > current.Price {
	//			top100ByPrice[index] = item
	//			current = item
	//		}
	//	}
	//}
	//
	//slices.SortFunc(top100ByPrice, func(i, j domain.Item) int {
	//	return cmp.Compare(i.Price, j.Price)
	//})
	//slices.Reverse(top100ByPrice)
	return top100ByPrice
}
