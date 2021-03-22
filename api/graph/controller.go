package graph

import (
	"context"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"

	"github.com/nasu/nasulog/api/graph/generated"
	"github.com/nasu/nasulog/api/infrastructure/dynamodb"
)

// Route sets echo route.
func Route(e *echo.Echo) {
	playgroundHandler := playground.Handler("GraphQL", "/gql/query")
	e.GET("/gql/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})

	ctx := context.TODO()
	db, err := dynamodb.GetDB(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	graphqlHandler := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{Resolvers: &Resolver{
				Ctx: ctx,
				DB:  db,
			}},
		),
	)
	e.POST("/gql/query", func(c echo.Context) error {
		graphqlHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}
