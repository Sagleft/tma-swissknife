package router

import (
	"errors"
	"log"
	"net/http"

	"github.com/Sagleft/tma-swissknife/rest"
	"github.com/gin-gonic/gin"
)

func onSuccess(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, rest.Success(data))
}

func onError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusInternalServerError, rest.ErrorMessage(err))
}

func maskError(ctx *gin.Context, err error, maskedErrInfo string) {
	log.Println(err)

	ctx.JSON(
		http.StatusInternalServerError,
		rest.ErrorMessage(errors.New(maskedErrInfo)),
	)
}
