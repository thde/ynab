# ynab

[![Go Reference](https://pkg.go.dev/badge/thde.io/ynab.svg)](https://pkg.go.dev/thde.io/ynab) [![test](https://github.com/thde/ynab/actions/workflows/test.yml/badge.svg)](https://github.com/thde/ynab/actions/workflows/test.yml) [![Go Report Card](https://goreportcard.com/badge/thde.io/ynab)](https://goreportcard.com/report/thde.io/ynab)

Go client for the [YNAB API](https://api.ynab.com/), generated from the official OpenAPI spec using [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen).

## Installation

```sh
go get thde.io/ynab
```

## Usage

```go
client, err := ynab.New(
    ynab.WithBearerToken("your-personal-access-token"),
)
if err != nil {
    log.Fatal(err)
}

resp, err := client.GetBudgetsWithResponse(ctx, nil)
```

Get a personal access token from your [YNAB developer settings](https://app.ynab.com/settings/developer).

Use `WithUserAgent` to identify your application in requests:

```go
client, err := ynab.New(
    ynab.WithBearerToken("your-personal-access-token"),
    ynab.WithUserAgent("my-app/1.0"),
)
```

## Rate limiting

The YNAB API [allows 200 requests per hour](https://api.ynab.com/#rate-limiting). Use `WithRateLimiter` to avoid hitting the limit:

```go
// Use the built-in default (200 req/hour, burst 1):
client, err := ynab.New(
    ynab.WithBearerToken("your-personal-access-token"),
    ynab.WithRateLimiter(nil),
)

// Or bring your own limiter:
limiter := rate.NewLimiter(rate.Every(time.Hour/200), 5)
client, err := ynab.New(
    ynab.WithBearerToken("your-personal-access-token"),
    ynab.WithRateLimiter(limiter),
)
```

## Delta requests

Most endpoints accept a `LastKnowledgeOfServer` parameter. When set, the server returns only entities that changed since that point, making repeated polling much cheaper. Each response includes a `ServerKnowledge` value to use on the next call.

```go
var knowledge int64

// First call, fetch everything.
resp, err := client.GetAccountsWithResponse(ctx, budgetID, &ynab.GetAccountsParams{})

// Persist resp.JSON200.Data.ServerKnowledge, then on subsequent calls:
resp, err = client.GetAccountsWithResponse(ctx, budgetID, &ynab.GetAccountsParams{
    LastKnowledgeOfServer: knowledge,
})
knowledge = resp.JSON200.Data.ServerKnowledge
```

See the [YNAB API docs](https://api.ynab.com/#deltas) for details.

## Development

```sh
make generate  # regenerate client from the YNAB OpenAPI spec
make test      # run tests
make lint      # run linters
make lint-fix  # run linters and auto-fix issues
make update    # update dependencies and regenerate
```

## License

MIT
