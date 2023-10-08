package middleware

import (
	"net/http"

	"go-bunrouter-gorm-example/infrastructure/httplib"
	"go-bunrouter-gorm-example/infrastructure/limiter"

	"github.com/uptrace/bunrouter"
)

func RateLimiterMiddleware(rateLimiter *limiter.RateLimiter) bunrouter.MiddlewareFunc {
	return func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			if rateLimiter.Allow() {
				return next(w, req)
			}
			return httplib.SetErrorResponse(w, http.StatusTooManyRequests, "rate limit exceeded")
		}
	}

}
