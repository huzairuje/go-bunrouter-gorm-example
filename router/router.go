package router

import (
	"net/http"

	"go-bunrouter-gorm-example/boot"
	"go-bunrouter-gorm-example/infrastructure/config"
	"go-bunrouter-gorm-example/infrastructure/httplib"
	"go-bunrouter-gorm-example/infrastructure/middleware"

	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

type HandlerRouter struct {
	Setup boot.HandlerSetup
}

func NewHandlerRouter(setup boot.HandlerSetup) InterfaceRouter {
	return &HandlerRouter{
		Setup: setup,
	}
}

type InterfaceRouter interface {
	RouterWithMiddleware() *bunrouter.Router
}

func notFoundHandler(w http.ResponseWriter, req bunrouter.Request) error {
	// render 404 custom response
	return httplib.SetErrorResponse(w, http.StatusNotFound, "Not Matching of Any Routes")
}

func methodNotAllowedHandler(w http.ResponseWriter, req bunrouter.Request) error {
	// render 404 custom response
	return httplib.SetErrorResponse(w, http.StatusMethodNotAllowed, "Method Not Allowed")
}

func (hr *HandlerRouter) RouterWithMiddleware() *bunrouter.Router {
	//add new instance for bun router and add not found handler
	//and method with not allowed handler
	c := bunrouter.New(
		bunrouter.WithNotFoundHandler(notFoundHandler),
		bunrouter.WithMethodNotAllowedHandler(methodNotAllowedHandler),
	)

	if config.Conf.LogMode {
		c.Use(reqlog.NewMiddleware(
			reqlog.WithEnabled(true),
			reqlog.WithVerbose(true),
			reqlog.FromEnv("BUNDEBUG"))).Verbose()
	}

	//grouping on root endpoint
	api := c.NewGroup("/api")

	api.Use(middleware.RateLimiterMiddleware(hr.Setup.Limiter))

	//grouping on "api/v1"
	v1 := api.NewGroup("/v1")

	//module health
	prefixHealth := v1.NewGroup("/health")
	hr.Setup.HealthHttp.GroupHealth(prefixHealth)

	//module article
	prefixArticle := v1.NewGroup("/articles")
	hr.Setup.ArticleHttp.GroupArticle(prefixArticle)

	return c

}
