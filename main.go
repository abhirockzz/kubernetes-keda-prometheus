package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var httpRequestsCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "http_requests",
	Help: "number of http requests",
})

const redisCounterName = "access_count"

var once sync.Once
var client *redis.Client

func init() {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		fmt.Printf("REDIS_HOST env var absent!")
		os.Exit(1)
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		fmt.Printf("REDIS_PORT env var absent!")
		os.Exit(1)
	}
	client = redis.NewClient(&redis.Options{Addr: host + ":" + port, PoolSize: 500})
	err := client.Ping().Err()
	if err != nil {
		fmt.Printf("Unable to connect to Redis at %s:%s", host, port)
		os.Exit(1)
	}
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		defer httpRequestsCounter.Inc()
		count, err := client.Incr(redisCounterName).Result()
		if err != nil {
			fmt.Println("Unable to increment redis counter", err)
			os.Exit(1)
		}
		resp := "Accessed on " + time.Now().String() + "\nAccess count " + strconv.Itoa(int(count))
		w.Write([]byte(resp))
	})

	http.ListenAndServe(":8080", nil)
}
