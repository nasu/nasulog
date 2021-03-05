package main

import (
	"github.com/nasu/nasulog/domain/article"
	"github.com/nasu/nasulog/domain/graphql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}: ${method} ${uri} status=${status}\n",
	}))

	article.Route(e)
	graphql.Route(e)
	e.Start(":8080")
}
