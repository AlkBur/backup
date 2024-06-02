package app

import (
	"backup/internal/config"
	"backup/internal/server"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	cfg := config.Server()

	r := server.Default()

	info := config.NewNotification()
	r.Notification(info)

	r.GET("/", func(c *server.Context) {
		c.String(http.StatusOK, "Hello Geektutu\n")
	})

	address := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Server on address %s\n", address)
	r.SendNotifications("Start backup", fmt.Sprintf("Server on address %s\n", address), "1")

	log.Fatal(r.Run(address))
}
