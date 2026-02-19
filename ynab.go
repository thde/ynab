// Package ynab provides a client for the YNAB (You Need A Budget) API.
package ynab

//go:generate go tool oapi-codegen -config oapi-codegen.yml https://api.ynab.com/papi/open_api_spec.yaml

import (
	"context"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const ProductionURL = ServerURLHTTPSAPIYNABComV1

const (
	headerAuthorization = "Authorization"
	headerUserAgent     = "User-Agent"
)

// WithBearerToken returns a [ClientOption] that sets the Authorization header with the provided token.
func WithBearerToken(token string) ClientOption {
	return WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set(headerAuthorization, "Bearer "+token)
		return nil
	})
}

// WithUserAgent returns a [ClientOption] that sets the User-Agent header with the provided token.
func WithUserAgent(agent string) ClientOption {
	return WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set(headerUserAgent, agent)
		return nil
	})
}

// defaultRateLimit is the YNAB API limit of 200 requests per hour.
// https://api.ynab.com/#rate-limiting
var defaultRateLimit = rate.Every(time.Hour / 200)

// WithRateLimiter returns a [ClientOption] that throttles outgoing requests using the provided
// [rate.Limiter]. This is useful for staying within the YNAB API limit of 200 requests per hour.
//
// If l is nil, a new limiter matching the YNAB rate limit (200 req/hour, burst 1) is created.
func WithRateLimiter(l *rate.Limiter) ClientOption {
	if l == nil {
		l = rate.NewLimiter(defaultRateLimit, 1)
	}

	return WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		return l.Wait(ctx)
	})
}

// New creates a new YNAB client pointing to the production endpoint with the provided options.
// In most cases, you should use this function to create a new client.
func New(opts ...ClientOption) (*ClientWithResponses, error) {
	return NewClientWithResponses(ProductionURL, opts...)
}
