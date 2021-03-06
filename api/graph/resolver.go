package graph

import (
	"context"

	"github.com/nasu/nasulog/api/infrastructure/dynamodb"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is resolver.
type Resolver struct {
	Ctx context.Context
	DB  *dynamodb.DB
}
