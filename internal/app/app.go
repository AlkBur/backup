package app

import (
	"backup/internal/config"
	"backup/internal/server"
	"fmt"
	"log"
	"net/http"
)

func Run() {
	cfg := config.Init()

	r := server.Default()
	r.GET("/", func(c *server.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})

	address := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server on address %s\n", address)
	log.Fatal(r.Run(address))
}
