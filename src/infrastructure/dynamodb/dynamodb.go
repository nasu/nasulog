package dynamodb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DB is a struct for database.
type DB struct {
	Client *dynamodb.Client
}

// GetDB gets DB struct.
func GetDB(ctx context.Context) (*DB, error) {
	resolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				URL: "http://localhost:9000",
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("ap-norteast-1"),
		config.WithEndpointResolver(resolver),
	)
	if err != nil {
		return nil, err
	}
	return &DB{dynamodb.NewFromConfig(cfg)}, nil
}

// SelectAll gets all data.
func (db DB) SelectAll(ctx context.Context, tableName string) ([]map[string]types.AttributeValue, error) {
	statement := fmt.Sprintf("SELECT * FROM %s", tableName)
	params := &dynamodb.ExecuteStatementInput{
		Statement: &statement,
	}
	res, err := db.Client.ExecuteStatement(ctx, params)
	if err != nil {
		return nil, err
	}
	//TODO: check NextToken
	return res.Items, nil
}

// UpsertMulti upserts for multi items in bulk.
func (db DB) UpsertMulti(ctx context.Context, tableName string, items []map[string]types.AttributeValue) error {
	for _, item := range items {
		err := db.upsert(ctx, tableName, item)
		if err != nil {
			return err
		}
	}
	return nil
}

// Insert inserts an item.
func (db DB) Insert(ctx context.Context, tableName string, item map[string]types.AttributeValue) error {
	conditionExpression := "attribute_not_exists(id)"
	params := &dynamodb.PutItemInput{
		TableName:           &tableName,
		Item:                item,
		ConditionExpression: &conditionExpression,
	}
	return db.putItem(ctx, params)
}

func (db DB) upsert(ctx context.Context, tableName string, item map[string]types.AttributeValue) error {
	params := &dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	}
	return db.putItem(ctx, params)
}

func (db DB) putItem(ctx context.Context, params *dynamodb.PutItemInput) error {
	_, err := db.Client.PutItem(ctx, params)
	if err != nil {
		return err
	}
	return nil
}
