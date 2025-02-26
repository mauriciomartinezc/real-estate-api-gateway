package middlewares

import (
	"github.com/juju/ratelimit"
	"log"
	"net/http"
)

// RateLimiterMiddleware limits the number of requests per second.
func RateLimiterMiddleware(reqsPerSecond int, burst int) func(http.Handler) http.Handler {
	if reqsPerSecond <= 0 || burst <= 0 {
		log.Printf("RateLimiterMiddleware: invalid rate limiter configuration (reqsPerSecond=%d, burst=%d)", reqsPerSecond, burst)
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Rate limiting not properly configured", http.StatusServiceUnavailable)
			})
		}
	}

	// Initialize the rate limiter
	bucket := ratelimit.NewBucketWithRate(float64(reqsPerSecond), int64(burst))
	if bucket == nil {
		log.Println("RateLimiterMiddleware: failed to initialize rate limiter bucket")
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		if next == nil {
			log.Println("RateLimiterMiddleware: next handler is nil")
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			})
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if bucket.TakeAvailable(1) == 0 {
				log.Printf("RateLimiterMiddleware: rate limit exceeded for %s %s", r.Method, r.URL.Path)
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			// Log successful request handling
			log.Printf("RateLimiterMiddleware: request allowed for %s %s", r.Method, r.URL.Path)

			// Call the next handler in the chain
			next.ServeHTTP(w, r)
		})
	}
}
