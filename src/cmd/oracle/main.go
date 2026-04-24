package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/bbombardella/arcana-oracle/internal/cache"
	"github.com/bbombardella/arcana-oracle/internal/handler"
	"github.com/bbombardella/arcana-oracle/internal/scaleway"
)

type app struct {
	card   *handler.CardHandler
	spread *handler.SpreadHandler
	astro  *handler.AstroHandler
}

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	ddb := dynamodb.NewFromConfig(cfg)
	cardCache := cache.NewDynamoDBCache(ddb, os.Getenv("DYNAMODB_TABLE"))
	scw := scaleway.NewClient(os.Getenv("SCW_API_URL"), os.Getenv("SCW_SECRET_KEY"))

	a := &app{
		card:   handler.NewCardHandler(scw, cardCache),
		spread: handler.NewSpreadHandler(scw),
		astro:  handler.NewAstroHandler(scw),
	}

	lambda.StartWithOptions(a.handle)
}

func (a *app) handle(ctx context.Context, event events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	req, err := buildRequest(ctx, event)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusInternalServerError}, err
	}

	rec := httptest.NewRecorder()

	switch event.RawPath {
	case "/oracle/card":
		a.card.ServeHTTP(rec, req)
	case "/oracle/spread":
		a.spread.ServeHTTP(rec, req)
	case "/oracle/astro":
		a.astro.ServeHTTP(rec, req)
	default:
		return events.APIGatewayV2HTTPResponse{StatusCode: http.StatusNotFound, Body: "not found"}, nil
	}

	headers := make(map[string]string, len(rec.Header()))
	for k, vs := range rec.Header() {
		headers[k] = strings.Join(vs, ",")
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: rec.Code,
		Headers:    headers,
		Body:       rec.Body.String(),
	}, nil
}

func buildRequest(ctx context.Context, event events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		event.RequestContext.HTTP.Method,
		event.RawPath,
		strings.NewReader(event.Body),
	)
	if err != nil {
		return nil, err
	}
	for k, v := range event.Headers {
		req.Header.Set(k, v)
	}
	return req, nil
}
