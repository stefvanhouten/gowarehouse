package main

import (
	"example/api/router"
	"example/api/types"
	"example/logger"
)

func main() {
	cfg := types.LoadEnv()
	logger.Setup(cfg)
	logger.DefaultLogger.Info("Starting server")

	router := router.New(cfg)
	router.Run("localhost:8080")
}
