package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {
	router := gin.Default()

	router.GET("/validators", GetValidators)

	router.Run()
}

func GetValidators(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"validators": []string{},
	})
}
