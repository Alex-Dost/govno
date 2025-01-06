package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// ClientLimiter хранит лимитаторы для каждого клиента
type ClientLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
	rate     rate.Limit
	burst    int
}

// NewClientLimiter создает новый ClientLimiter
func NewClientLimiter(r rate.Limit, b int) *ClientLimiter {
	return &ClientLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     r,
		burst:    b,
	}
}

// GetLimiter возвращает лимитер для конкретного IP
func (c *ClientLimiter) GetLimiter(ip string) *rate.Limiter {
	c.mu.Lock()
	defer c.mu.Unlock()

	limiter, exists := c.limiters[ip]
	if !exists {
		limiter = rate.NewLimiter(c.rate, c.burst)
		c.limiters[ip] = limiter
	}

	return limiter
}

func main() {
	clientLimiter := NewClientLimiter(5, 10) // 5 запросов в секунду, максимум 10 в буфере

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)

		limiter := clientLimiter.GetLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		fmt.Fprintf(w, "Welcome! Your IP is %s", ip)
	})

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server is running on port 8080...")
	log.Fatal(server.ListenAndServe())
}

// getIP получает IP адрес клиента
func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For") // Для прокси-серверов
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}
