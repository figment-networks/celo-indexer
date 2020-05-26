package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/figment-networks/celo-indexer/grpc/validator"
)

func Run() {
	router := gin.Default()

	router.GET("/validators", GetValidators)

	router.Run()
}

const proxyUrl = "localhost:50051"

func GetValidators(ctx *gin.Context) {
	conn, err := grpc.Dial(proxyUrl, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic("could not connect to gRPC proxy")
	}
	defer conn.Close()

	client := validator.NewValidatorServiceClient(conn)

	response, err := client.GetByHeight(context.Background(),
		&validator.GetByHeightRequest{Height: 0})

	ctx.JSON(http.StatusOK, response)
}
