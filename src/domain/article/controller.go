package article

import (
	"context"
	"log"
	"net/http"

	"github.com/nasu/nasulog/infrastructure/dynamodb"

	"github.com/labstack/echo/v4"
)

// Route set echo route.
func Route(e *echo.Echo) {
	e.GET("/article", all)
	e.GET("/article/:id", one)
	e.POST("/article", create)
}

//TODO: write API reference
//TODO: write e2e test

func one(c echo.Context) error {
	ctx := context.TODO()
	id := c.Param("id")

	client, err := dynamodb.GetClient(ctx)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	article, err := selectOne(ctx, client, id)
	if err != nil {
		log.Printf("failed to get article, %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, article)
}

func all(c echo.Context) error {
	ctx := context.TODO()
	client, err := dynamodb.GetClient(ctx)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	articles, err := selectAll(ctx, client)
	if err != nil {
		log.Printf("failed to get articles, %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, articles)
}

func create(c echo.Context) error {
	ctx := context.TODO()

	var article Article
	if err := c.Bind(&article); err != nil {
		log.Printf("unable to bind body, %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	client, err := dynamodb.GetClient(ctx)
	if err != nil {
		log.Printf("unable to load SDK config, %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}

	insertedArticle, err := insert(ctx, client, &article)
	if err != nil {
		log.Printf("failed to insert an article, %v", err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusCreated, insertedArticle)
}
