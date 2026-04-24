package cache

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBCache caches card interpretations.
// PK format: "card#<cardName>#<reversed>#<lang>"  e.g. "card#Le Bateleur#false#fr"
type DynamoDBCache struct {
	client *dynamodb.Client
	table  string
}

func NewDynamoDBCache(client *dynamodb.Client, table string) *DynamoDBCache {
	return &DynamoDBCache{client: client, table: table}
}

type cacheItem struct {
	PK       string `dynamodbav:"pk"`
	Response string `dynamodbav:"response"`
}

func buildKey(cardID string, reversed bool, lang string) string {
	r := "false"
	if reversed {
		r = "true"
	}
	return fmt.Sprintf("card#%s#%s#%s", cardID, r, lang)
}

func (c *DynamoDBCache) Get(ctx context.Context, cardID string, reversed bool, lang string) (string, bool, error) {
	pk := buildKey(cardID, reversed, lang)

	out, err := c.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(c.table),
		Key: map[string]dbtypes.AttributeValue{
			"pk": &dbtypes.AttributeValueMemberS{Value: pk},
		},
	})
	if err != nil {
		return "", false, err
	}
	if out.Item == nil {
		return "", false, nil
	}

	var it cacheItem
	if err := attributevalue.UnmarshalMap(out.Item, &it); err != nil {
		return "", false, err
	}
	return it.Response, true, nil
}

func (c *DynamoDBCache) Set(ctx context.Context, cardID string, reversed bool, lang, response string) error {
	it := cacheItem{
		PK:       buildKey(cardID, reversed, lang),
		Response: response,
	}
	av, err := attributevalue.MarshalMap(it)
	if err != nil {
		return err
	}
	_, err = c.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(c.table),
		Item:      av,
	})
	return err
}
