package middleware

import (
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	rate     int
	burst    int
	window   time.Duration
}

type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewRateLimiter(rate int, burst int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*Visitor),
		rate:     rate,
		burst:    burst,
		window:   window,
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rl.mu.Lock()

		v, exists := rl.visitors[ip]
		if !exists {
			limiter := rate.NewLimiter(rate.Every(rl.window/time.Duration(rl.rate)), rl.burst)
			rl.visitors[ip] = &Visitor{limiter, time.Now()}
			v = rl.visitors[ip]
		}
		v.lastSeen = time.Now()
		rl.mu.Unlock()

		if !v.limiter.Allow() {
			http.Error(w, "Rate limit Exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
