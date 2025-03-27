package mw

import (
	"net/http"
	"sync"
	"time"
)

type window struct {
    count   int
    resetAt time.Time
}

var (
    requests map[string]*window
    mu       sync.RWMutex
    maxReq   int
    duration time.Duration
)

func SetRateLimit(max int, d time.Duration) {
    requests = make(map[string]*window)
    maxReq = max
    duration = d
}

func RateLimitMiddleware(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        now := time.Now()

        mu.Lock()
        if win, exists := requests[ip]; exists {
            if now.After(win.resetAt) {
                requests[ip] = &window{
                    count:   1,
                    resetAt: now.Add(duration),
                }
            } else {
                win.count++
            }
        } else {
            requests[ip] = &window{
                count:   1,
                resetAt: now.Add(duration),
            }
        }

        win := requests[ip]
        remaining := maxReq - win.count
        mu.Unlock()

        if remaining < 0 {
            http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        handler.ServeHTTP(w, r)
    })
}
