package routes

import (
	"encoding/json"
	"io"
	"ipbans/blocklist"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CheckIpBody struct {
	Ip string
}

func POSTCheck(ctx *gin.Context) {
	rawBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to read body",
			"error":   err,
		})
		return
	}

	var body CheckIpBody
	err = json.Unmarshal(rawBody, &body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to parse json",
			"error":   err,
		})
		return
	}

	ip := net.ParseIP(body.Ip)
	if ip == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid IP address",
		})
		return
	}

	// Check if ip is in the array
	isBlocked := blocklist.Contains(ip)

	ctx.JSON(http.StatusOK, gin.H{
		"ip":        body.Ip,
		"isBlocked": isBlocked,
	})
}
