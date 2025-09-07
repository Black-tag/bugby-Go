package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/blacktag/bugby-Go/internal/caching"
	
)

func CachingMiddleware(expiration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := slog.Default().With("middleware", "caching")
			cacheKey := r.Method + ":" + r.URL.String()
			logger = logger.With("cacheKey", cacheKey)

			// 1) Try cache
			cached, err := caching.GetCached(cacheKey)
			if err != nil {
				logger.Error("cache error", "err", err)
			}
			if cached != nil {
				logger.Info("cache hit")
				for k, v := range cached.Headers {
					w.Header().Set(k, v)
				}
				w.WriteHeader(cached.StatusCode)
				_, _ = w.Write(cached.Body)
				return
			}

			// 2) Miss: record and forward
			logger.Info("cache miss")
			rec := caching.NewRecorder(w)
			next.ServeHTTP(rec, r)

			// 3) Store
			resp := caching.ResponseWrapper(rec.StatusCode, rec.Header(), rec.Body)
			if err := resp.SetCached(cacheKey, 5*time.Minute); err != nil {
				logger.Error("failed to set cache", "err", err)
			}
		})
	}
}
