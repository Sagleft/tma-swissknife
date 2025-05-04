package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Sagleft/tma-swissknife/rest"
	"github.com/gin-gonic/gin"
)

var errNotFound = errors.New("not found")

type Router interface {
	// NOTE: it's blocking method
	Serve()
}

type RouterHandler interface {
	HandleRequest(rest.Request) (rest.Message, error)
}

type router struct {
	engine *gin.Engine
}

func New() *router {
	return &router{
		engine: gin.New(),
	}
}

func (r *router) Serve(host, port string, h RouterHandler) error {
	r.engine.GET("/ping", func(ctx *gin.Context) {
		onSuccess(ctx, "pong")
	})

	r.engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, rest.ErrorMessage(errNotFound))
	})

	r.engine.POST("/api", func(ctx *gin.Context) {
		var req rest.Request
		if err := ctx.ShouldBind(&req); err != nil {
			maskError(ctx, fmt.Errorf("parse req: %w", err), "failed to process request")
			return
		}

		response, err := h.HandleRequest(req)
		if err != nil {
			maskError(ctx, fmt.Errorf("%q: %w", req.Method, err), "failed to handle request")
			return
		}

		ctx.JSON(http.StatusOK, response)
	})

	return r.engine.Run(fmt.Sprintf("%s:%s", host, port))
}
