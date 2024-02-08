package integrations

import (
	"context"
	"reliability_demo_app/internal/domain"
	"reliability_demo_app/internal/httpclient"
	"reliability_demo_app/internal/usecases"

	"github.com/melisource/fury_go-core/pkg/rusty"
)

var _ usecases.BusinessFlowProcesor[domain.PriceChangeRequest] = &ChangePriceFlowProcessor{}

func NewChangePriceFlowProcessor(factory httpclient.EndpointFactory) *ChangePriceFlowProcessor {
	return &ChangePriceFlowProcessor{
		endpoint: factory.Build("/v1/price-changes-aproval-flow"),
	}
}

type ChangePriceFlowProcessor struct {
	endpoint *rusty.Endpoint
}

func (p *ChangePriceFlowProcessor) BeginFlow(ctx context.Context, input domain.PriceChangeRequest, callbackURL string) error {
	body := PriceChangeRequestMessage{
		TransactionID: input.TransactionID,
		ItemID:        input.ItemID,
		Price:         input.Price,
		CallbackUrl:   callbackURL,
	}
	res, err := p.endpoint.Post(ctx, rusty.WithBody(body))
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return err
	}

	return nil
}

type PriceChangeRequestMessage struct {
	TransactionID string  `json:"transaction_id"`
	ItemID        string  `json:"item_id"`
	Price         float64 `json:"price"`
	CallbackUrl   string  `json:"callback_url"`
}
