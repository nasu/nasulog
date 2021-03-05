package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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
