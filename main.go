package main

import (
	"ipbans/blocklist"
	"ipbans/routes"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	slog.Info("Hello world!")

	quitSignal := make(chan struct{})
	go blocklist.ScheduleRefreshBlocklists(quitSignal)

	blocklist.RefreshBlocklists()
	if !blocklist.BlocklistReady {
		slog.Warn("Blocklist is not ready")
	}

	address, exist := os.LookupEnv("ADDR")
	if !exist {
		address = "0.0.0.0:3000"
	}

	router := gin.Default()

	router.GET("/", routes.GETIndex)
	router.POST("/check", routes.POSTCheck)

	router.Run(address)
}
