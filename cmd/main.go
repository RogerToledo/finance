package main

import (
	"github.com/me/application/go/finance/internal/config"
	"github.com/me/application/go/finance/internal/config/logger"
	"github.com/me/application/go/finance/internal/handler"
)

func main() {
	logger.InitLogger()
	
	if err := config.Load(); err != nil {
		panic(err)
	}

	router.InitializeHandler()
}``