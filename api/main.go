package main

import (
	"log"
	"os"

	"github.com/nasu/nasulog/api/graph"
	"github.com/nasu/nasulog/api/infrastructure/dynamodb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	url, _ := os.LookupEnv("DYNAMODB_URL")
	region, _ := os.LookupEnv("DYNAMODB_REGION")
	dynamodb.InjectConstant(url, region)
	log.Println("DYNAMODB_URL:" + dynamodb.DYNAMODB_URL)
	log.Println("DYNAMODB_REGION:" + dynamodb.DYNAMODB_REGION)
}

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}: ${method} ${uri} status=${status}\n",
	}))

	graph.Route(e)
	e.Start(":8080")
}
