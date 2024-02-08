package main

import (
	"fmt"
	"os"
	"reliability_demo_app/internal/http"
	"reliability_demo_app/internal/httpclient"
	"reliability_demo_app/internal/integrations"
	"reliability_demo_app/internal/usecases"

	"github.com/melisource/fury_go-core/pkg/log"
	"github.com/melisource/fury_go-platform/pkg/fury"
)

func main() {
	level := log.DebugLevel
	opts := []fury.AppOptFunc{}
	opts = append(opts, fury.WithLogLevel(level))
	app, err := fury.NewWebApplication(opts...)
	if err != nil {
		panic(err)
	}

	backendURL := "http://drs-product-catalog-backend.melisystems.com"
	myselfURL := fmt.Sprintf("http://%s-%s.melisystems.com", os.Getenv("SCOPE"), os.Getenv("APPLICATION"))
	if os.Getenv("ENV") == "local" {
		backendURL = "http://localhost:10000"
		myselfURL = "http://localhost:8080"
	}

	endpointFactory := httpclient.NewEndpointFactory(backendURL)
	searchEngine := integrations.NewSearchEngine(endpointFactory)
	itemApi := integrations.NewItemsApi(endpointFactory)
	searchService := usecases.NewSearchService(searchEngine)
	changePriceFlowProcessor := integrations.NewChangePriceFlowProcessor(endpointFactory)
	itemService := usecases.NewItemService(
		itemApi,
		itemApi,
		changePriceFlowProcessor,
		myselfURL,
	)
	reportService := usecases.NewReportService(searchService)

	controllers := []http.Controller{
		http.NewItemController(itemService),
		http.NewSearchController(searchService),
		http.NewReportController(reportService),
	}

	for _, controller := range controllers {
		controller.AddRoutes(app.Router)
	}

	if err := app.Run(); err != nil {
		panic(err)
	}
}
