package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/grassrootseconomics/ussd-canary-proxy/internal/api"
	"github.com/uptrace/bunrouter"
)

var (
	ErrNotAuthorized = errors.New("not authorized")
)

func initRouter(api *api.API) *bunrouter.Router {
	router := bunrouter.New()

	if ko.Bool("service.metrics") {
		metrics := router.NewGroup("/metrics").Compat()
		metrics.GET("/", api.MetricsHandler())
	}

	cacheAPI := router.NewGroup("/api/cache")
	cacheAPI.Use(errorHandler).Use(authMiddleware)
	cacheAPI.POST("/add", api.AddToCacheHandler)
	cacheAPI.GET("/get/:phoneNumber", api.GetFromCacheHandler)

	ussdIngress := router.NewGroup("/ussd").Compat()
	ussdIngress.POST(fmt.Sprintf("/%s", ko.MustString("webhook.secret")), api.USSDProxyHandler())

	return router
}

func authMiddleware(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		if req.Header.Get("x-api-key") != ko.MustString("service.api_key") {
			return ErrNotAuthorized
		}

		return next(w, req)
	}
}

func errorHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		err := next(w, req)

		switch err {
		case nil:
		case ErrNotAuthorized:
			w.WriteHeader(http.StatusTooManyRequests)
			_ = bunrouter.JSON(w, bunrouter.H{
				"message": "Not Authorized",
			})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			_ = bunrouter.JSON(w, bunrouter.H{
				"message": err.Error(),
			})
		}

		return err
	}
}
