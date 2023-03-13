package main

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/rafli-lutfi/go-commerce/config"
	"github.com/rafli-lutfi/go-commerce/routes"
)

func init() {
	config.LoadEnv()
	config.LoadDatabase()
}

func main() {
	db := config.GetDBConnection()
	r := gin.Default()

	wg := sync.WaitGroup{}

	wg.Add(1)

	go func() {
		defer wg.Done()
		routes.RunServer(db, r)

	}()

	r.Run()
}
