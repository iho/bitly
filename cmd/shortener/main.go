package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iho/bitly/conf"
	"github.com/iho/bitly/handlers"
	"github.com/iho/bitly/shortener"
	"github.com/iho/bitly/storage"
	"go.uber.org/fx"
)

func main() {
	fmt.Println("Starting...")
	fx.New(
		storage.Module,
		shortener.Module,
		handlers.Module,
		conf.Module,
		fx.Invoke(func(r *gin.Engine, c *conf.Config) { r.Run(c.HostName) }),
	).Run()

}
