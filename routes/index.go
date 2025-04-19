package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GETIndex(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the IP Bans API! Use /check to check if an IP is blocked.",
	})
}
