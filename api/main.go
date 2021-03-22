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
	if val, ok := os.LookupEnv("DYNAMODB_URL"); !ok {
		log.Fatalln("require ENV:DYNAMODB_URL")
	} else {
		dynamodb.InjectEndpointURL(val)
		log.Println("DYNAMODB_URL:" + val)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339}: ${method} ${uri} status=${status}\n",
	}))

	graph.Route(e)
	e.Start(":8080")
}
