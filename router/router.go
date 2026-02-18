package router

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"text/template"

	"github.com/Sagleft/tma-swissknife/rest"
	"github.com/gin-gonic/gin"
)

const (
	HttpGET    HttpMethod = "GET"
	HttpPOST   HttpMethod = "POST"
	HttpPUT    HttpMethod = "PUT"
	HttpDELETE HttpMethod = "DELETE"
)

type HttpMethod string

var errNotFound = errors.New("not found")

// all setups must be called before Serve()
type Router interface {
	SetupRoutes(routes []Route)
	SetupTemplates(cfg TemplateConfig) error
	SetupTLS(TLSConfig) error
	SetupErrorHandler(func(error))
	SetTrustedProxies([]string) error

	// NOTE: it's blocking method
	Serve(host, port string, h RouterHandler) error
}

type RouterHandler interface {
	HandleRequest(rest.Request) (rest.Message, error)
}

type router struct {
	engine     *gin.Engine
	errHandler func(error)
	tls        TLSConfig
}

type TLSConfig struct {
	Enabled      bool   `json:"enabled"`
	CertFilepath string `json:"certFilepath"` // example: localhost.pem
	KeyFilepath  string `json:"keyFilepath"`  // example: localhost-key.pem
}

func New() Router {
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
			handleError(ctx, fmt.Errorf("parse req: %w", err))
			return
		}

		response, err := h.HandleRequest(req)
		if err != nil {
			handleError(ctx, fmt.Errorf("%q: %w", req.Method, err))
			return
		}

		ctx.JSON(http.StatusOK, response)
	})

	addr := fmt.Sprintf("%s:%s", host, port)

	if r.tls.Enabled {
		return r.engine.RunTLS(addr, r.tls.CertFilepath, r.tls.KeyFilepath)
	}

	return r.engine.Run(addr)
}

type Route struct {
	// required
	Endpoint string
	Handler  gin.HandlerFunc

	// optional
	HttpMethod HttpMethod
}

func (r *router) SetupErrorHandler(f func(error)) {
	r.errHandler = f
	r.engine.Use(r.errorMiddleware)
}

func (r *router) errorMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Получаем стек вызовов
			stack := make([]byte, 4096)
			n := runtime.Stack(stack, true)

			if r.errHandler != nil {
				r.errHandler(fmt.Errorf("panic: %v, stack: %s", err, stack[:n]))
			}
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	c.Next()

	for _, e := range c.Errors {
		if r.errHandler != nil {
			r.errHandler(e.Err)
		}
	}

	if len(c.Errors) > 0 && !c.IsAborted() {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (r *router) SetupRoutes(routes []Route) {
	for _, routeData := range routes {
		r.engine.Handle(
			ternary(
				routeData.HttpMethod == "",
				string(HttpGET),
				string(routeData.HttpMethod),
			),
			routeData.Endpoint,
			routeData.Handler,
		)
	}
}

type TemplateConfig struct {
	Path             string   `json:"path"`            // example: "templates/*"
	CustomDelimeters []string `json:"delimeters"`      // example: "{[{", "}]}"
	StaticAssetsPath string   `json:"staticFilesPath"` // example: "./public"
}

func (r *router) SetupTemplates(cfg TemplateConfig) error {
	if len(cfg.CustomDelimeters) > 1 {
		r.engine.Delims(cfg.CustomDelimeters[0], cfg.CustomDelimeters[1])
	}

	if cfg.StaticAssetsPath != "" {
		if err := r.registerStaticAssets(cfg.StaticAssetsPath); err != nil {
			return fmt.Errorf("reg assets: %w", err)
		}
	}

	r.engine.LoadHTMLGlob(cfg.Path)
	return nil
}

func (r *router) registerStaticAssets(assetsPath string) error {
	hashes, err := HashAssets(assetsPath)
	if err != nil {
		return fmt.Errorf("hash assets: %w", err)
	}

	// asset version system
	r.engine.SetFuncMap(template.FuncMap{
		"versioned": func(path string) string {
			if hash, isExists := hashes[path]; isExists {
				return "/assets/" + path + "?v=" + hash
			}
			return "/assets/" + path
		},
		"getIcon": func(tag string) string {
			if tag == "" {
				return "star"
			}
			return tag
		},
	})

	r.engine.Static("/assets", assetsPath)
	return nil
}

func (r *router) SetupTLS(cfg TLSConfig) error {
	if cfg.CertFilepath == "" {
		return errors.New("cert file path not set")
	}
	if cfg.KeyFilepath == "" {
		return errors.New("cert key file path not set")
	}

	r.tls = cfg
	return nil
}

func (r *router) SetTrustedProxies(p []string) error {
	return r.engine.SetTrustedProxies(p)
}
