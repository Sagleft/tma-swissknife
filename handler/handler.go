package handler

import (
	"fmt"

	"github.com/Sagleft/tma-swissknife/rest"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Init(methods map[string]HandlerFunc)

	HandleRequest(*gin.Context, rest.Request) (rest.Message, error)
}

type handler struct {
	methods map[string]HandlerFunc
}

// data - data from json request
type HandlerFunc func(ctx *gin.Context, data map[string]any) (rest.Message, error)

func New() Handler {
	return &handler{methods: make(map[string]HandlerFunc)}
}

func (h *handler) Init(methods map[string]HandlerFunc) {
	h.methods = methods
}

func (h *handler) HandleRequest(
	ctx *gin.Context,
	req rest.Request,
) (rest.Message, error) {
	method, ok := h.methods[req.Method]
	if !ok {
		return rest.Message{}, fmt.Errorf("method %q not found", req.Method)
	}
	return method(ctx, req.Data)
}
