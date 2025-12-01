package caching

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

type CachedResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func InitCache() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return errors.New("REDISURL env not set")
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}
	Client = redis.NewClient(opt)

	_, err = Client.Ping(Ctx).Result()
	return err

}

// points for caching
// make a request
// check if the cache exists
// if its in cache return cached response
// if its not in cache do req catch req set it in cache

func GetCached(cacheKey string) (*CachedResponse, error) {
	logger := slog.Default().With(
		"caching function", "retrieving from Cache",
	)
	logger.Info("retrieving response from Cache")

	val, err := Client.Get(Ctx, cacheKey).Result()
	if err == redis.Nil {

		return nil, err
	}
	if err != nil {

		return nil, err
	}
	var cached CachedResponse
	err = json.Unmarshal([]byte(val), &cached)
	if err != nil {
		err = errors.New("couldnt unmarshal json")
		return nil, err
	}

	return &cached, nil

}

func (cache *CachedResponse) SetCached(cacheKey string, expiration time.Duration) error {
	logger := slog.Default().With(
		"caching function", "setting response into cache with cacheKey",
	)
	logger.Info("Setting  response to Cache")
	jsonData, err := json.Marshal(cache)
	if err != nil {
		logger.Error("couldnt marshal Cached response into json")
		logger.Error("couldnt marshal json", "err", err)
		return err

	}
	cachedresult := Client.Set(Ctx, cacheKey, jsonData, expiration).Err()
	if cachedresult != nil {
		logger.Error("couldnt fetch cached result")
		logger.Error("couldnt set cache in redis", "err", err)
		return nil
	}
	return cachedresult

}

func ResponseWrapper(statuscode int, headers http.Header, body []byte) *CachedResponse {

	headerMap := make(map[string]string)
	for k, v := range headers {
		if len(v) > 0 {
			headerMap[k] = v[0]
		}

	}
	return &CachedResponse{
		StatusCode: statuscode,
		Headers:    headerMap,
		Body:       body,
	}

}

type Recorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func NewRecorder(w http.ResponseWriter) *Recorder {
	return &Recorder{ResponseWriter: w, StatusCode: http.StatusOK}
}

func (r *Recorder) WriteHeader(code int) {
	r.StatusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *Recorder) Write(b []byte) (int, error) {
	r.Body = append(r.Body, b...)
	return r.ResponseWriter.Write(b)

}
