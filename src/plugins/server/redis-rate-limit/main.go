// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

// pluginName is the plugin name
var pluginName = "krakend-server-example"

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = registerer(pluginName)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(ctx context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	config, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return h, errors.New("configuration not found")
	}

	redisHost, _ := config["redis_host"].(string)

	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})
	_ = rdb.FlushDB(ctx).Err()
	limiter := redis_rate.NewLimiter(rdb)

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		apiKey := req.Header.Get("X-API-KEY")

		if len(apiKey) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Please set X-API-KEY into HTTP Headers."))
			return
		}

		rateLimitRes, err := limiter.Allow(ctx, apiKey, redis_rate.PerMinute(2))

		if err != nil {
			panic(err)
		}

		if rateLimitRes.Allowed <= 0 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Too many requests."))
			return
		}

		h.ServeHTTP(w, req)
	}), nil
}

func main() {}

// This logger is replaced by the RegisterLogger method to load the one from KrakenD
var logger Logger = noopLogger{}

func (registerer) RegisterLogger(v interface{}) {
	l, ok := v.(Logger)
	if !ok {
		return
	}
	logger = l
	logger.Debug(fmt.Sprintf("[PLUGIN: %s] Logger loaded", HandlerRegisterer))
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Critical(v ...interface{})
	Fatal(v ...interface{})
}

// Empty logger implementation
type noopLogger struct{}

func (n noopLogger) Debug(_ ...interface{})    {}
func (n noopLogger) Info(_ ...interface{})     {}
func (n noopLogger) Warning(_ ...interface{})  {}
func (n noopLogger) Error(_ ...interface{})    {}
func (n noopLogger) Critical(_ ...interface{}) {}
func (n noopLogger) Fatal(_ ...interface{})    {}
