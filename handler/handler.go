package handler

import (
	"fmt"

	"github.com/Sagleft/tma-swissknife/rest"
)

type Handler interface {
	Init(methods map[string]HandlerFunc)

	HandleRequest(rest.Request) (rest.Message, error)
}

type handler struct {
	methods map[string]HandlerFunc
}

type HandlerFunc func(data any) (rest.Message, error)

func New() (*handler, error) {
	h := &handler{}
	return h, nil
}

func (h *handler) Init(methods map[string]HandlerFunc) {
	h.methods = methods
}

func (h *handler) HandleRequest(req rest.Request) (rest.Message, error) {
	method, ok := h.methods[req.Method]
	if !ok {
		return rest.Message{}, fmt.Errorf("method %q not found", req.Method)
	}
	return method(req.Data)
}
