package dynamodb

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DB is a struct for database.
type DB struct {
	Client *dynamodb.Client
}

var (
	DYNAMODB_URL string
)

func InjectEndpointURL(url string) {
	DYNAMODB_URL = url
}

// GetDB gets DB struct.
func GetDB(ctx context.Context) (*DB, error) {
	if DYNAMODB_URL == "" {
		return nil, fmt.Errorf("Not found dynamodb endpoint.")
	}
	resolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if service == dynamodb.ServiceID {
			return aws.Endpoint{
				URL: DYNAMODB_URL,
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

// SelectByPK gets one data with primary keys.
func (db DB) SelectByPK(ctx context.Context, tableName string, key map[string]types.AttributeValue, consistentRead bool) (map[string]types.AttributeValue, error) {
	res, err := db.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:      &tableName,
		Key:            key,
		ConsistentRead: &consistentRead,
	})
	if err != nil {
		return nil, err
	}
	return res.Item, nil
}

// SelectAll gets all data.
func (db DB) SelectAll(ctx context.Context, tableName, partitionKey string) ([]map[string]types.AttributeValue, error) {
	statement := fmt.Sprintf("SELECT * FROM %s WHERE partition_key=?", tableName)
	params := []types.AttributeValue{&types.AttributeValueMemberS{Value: partitionKey}}
	input := &dynamodb.ExecuteStatementInput{
		Statement:  &statement,
		Parameters: params,
	}
	res, err := db.Client.ExecuteStatement(ctx, input)
	if err != nil {
		return nil, err
	}
	//TODO: check NextToken
	return res.Items, nil
}

// SelectBySortKeys gets data with partition_key and multi sort_key
func (db DB) SelectBySortKeys(ctx context.Context, tableName, partitionKey string, sortKeys []string) ([]map[string]types.AttributeValue, error) {
	if len(sortKeys) == 0 {
		return []map[string]types.AttributeValue{}, nil
	}

	params := []types.AttributeValue{&types.AttributeValueMemberS{Value: partitionKey}}
	placeHolders := make([]string, len(sortKeys), len(sortKeys))
	for i, key := range sortKeys {
		placeHolders[i] = "?"
		params = append(params, &types.AttributeValueMemberS{Value: key})
	}
	sortKeyPlaceHolder := strings.Join(placeHolders, ",")
	statement := fmt.Sprintf("SELECT * FROM %s WHERE partition_key=? AND sort_key IN [%s]", tableName, sortKeyPlaceHolder)

	input := &dynamodb.ExecuteStatementInput{
		Statement:  &statement,
		Parameters: params,
	}
	res, err := db.Client.ExecuteStatement(ctx, input)
	if err != nil {
		return nil, err
	}
	//TODO: check NextToken
	return res.Items, nil
}

// UpsertMulti upserts for multi items in bulk.
func (db DB) UpsertMulti(ctx context.Context, tableName string, items []map[string]types.AttributeValue) error {
	for _, item := range items {
		params := &dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      item,
		}
		err := db.putItem(ctx, params)
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

// DeleteByPK deletes an item with pk.
func (db DB) DeleteByPK(ctx context.Context, tableName string, key map[string]types.AttributeValue) error {
	_, err := db.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key:       key,
	})
	if err != nil {
		return err
	}
	return nil
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
