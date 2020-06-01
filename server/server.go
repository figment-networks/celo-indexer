package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/figment-networks/celo-indexer/client"
)

func Run() {
	router := gin.Default()

	router.GET("/validators", GetValidators)

	router.Run()
}

func GetValidators(ctx *gin.Context) {
	client, err := client.New("localhost:50051")
	if err != nil {
		panic(err)
	}
	defer client.Close()

	validators, err := client.Validator.GetByHeight(0)
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, validators)
}
