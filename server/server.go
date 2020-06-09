package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/figment-networks/celo-indexer/model"
	"github.com/figment-networks/celo-indexer/store"
)

func Run(store *store.Store) {
	router := gin.Default()

	router.GET("/validators", func(ctx *gin.Context) {
		var validators []model.Validator

		store.Db.Find(&validators)

		ctx.JSON(http.StatusOK, validators)
	})

	router.Run()
}
