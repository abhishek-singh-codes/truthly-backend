package tail38

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/redis/go-redis/v9"
)

type GeoClient interface {
	SetPoint(
		ctx context.Context,
		collection string,
		id string,
		lat float64,
		long float64,
	) error

	// NearByImages(
	// 	ctx context.Context,
	// 	collection string,
	// 	long float32,
	// 	lat float32,
	// 	rangeUnit int,
	// 	limit int,
	// 	cursor int,
	// ) ([]string, int, bool, error)
}

// this time I am using redis based geo client
type geoClient struct {
	rdb    *redis.Client
	logger *slog.Logger
}

var (
	once     sync.Once
	instance GeoClient
)

// this func will initialize/ return the GeoClient interface
func InitGeoClient(
	host string,
	port int,
	logger *slog.Logger,
) GeoClient {
	once.Do(func() {
		addr := fmt.Sprintf("%s:%d", host, port)

		rdb := redis.NewClient(&redis.Options{
			Addr:         addr,
			PoolSize:     30,
			MinIdleConns: 10,
		})

		logger.Warn(
			"GeoClient connecting",
			"addr", addr,
		)

		ctx := context.Background()
		if err := rdb.Ping(ctx).Err(); err != nil {
			logger.Error("Redis connection failed", "error", err)
			panic(err)
		}

		instance = &geoClient{
			rdb:    rdb,
			logger: logger,
		}

		logger.Info("Tail38 Initialized", "addr", addr)
	})

	return instance
}
