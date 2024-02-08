package httpclient

import (
	"time"

	"github.com/melisource/fury_go-core/pkg/rusty"
	"github.com/melisource/fury_go-core/pkg/transport"
	"github.com/melisource/fury_go-core/pkg/transport/httpclient"
)

const (
	retryMax                  = 4
	dialTimeOut               = 1 * time.Second
	requestTimeOut            = 30 * time.Second
	defaultPoolName           = "reliability_demo"
	DefaultBackoffMinDuration = 250 * time.Millisecond
	DefaultBackoffMaxDuration = 5 * time.Second
)

type EndpointFactory interface {
	Build(pattern string) *rusty.Endpoint
}

var _ EndpointFactory = &RustyEndpointFactory{}

func NewEndpointFactory(
	baseURL string,
) *RustyEndpointFactory {
	return &RustyEndpointFactory{
		baseURL: baseURL,
	}
}

type RustyEndpointFactory struct {
	baseURL string
}

func (f *RustyEndpointFactory) Build(pattern string) *rusty.Endpoint {
	transportOpts := []transport.Option{}

	clientOpts := []httpclient.OptionRetryable{
		httpclient.WithTimeout(requestTimeOut),
		httpclient.WithTransport(transport.NewPooled(defaultPoolName, transportOpts...)),
		httpclient.WithBackoffStrategy(httpclient.ExponentialBackoff(DefaultBackoffMinDuration, DefaultBackoffMaxDuration)),
	}

	httpClient := httpclient.NewRetryable(retryMax, clientOpts...)

	opts := []rusty.EndpointOption{
		rusty.WithHeader("Content-Type", "application/json"),
	}

	endpoint, err := rusty.NewEndpoint(httpClient, rusty.URL(f.baseURL, pattern), opts...)
	if err != nil {
		panic(err)
	}

	return endpoint
}
