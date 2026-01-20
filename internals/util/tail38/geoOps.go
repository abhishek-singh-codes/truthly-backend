package tail38

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// set point functionality
func (t *geoClient) SetPoint(
	ctx context.Context,
	collection string,
	id string,
	long float64,
	lat float64,
) error {

	err := t.rdb.GeoAdd(ctx, collection, &redis.GeoLocation{
		Name:      id,
		Longitude: long,
		Latitude:  lat,
	}).Err()

	t.logger.Info(
		"Inserted image location into Redis GEO",
		"imageId", id,
		"lat", lat,
		"long", long,
		"collection", collection,
	)

	if err != nil {
		t.logger.Error("Redis GEOADD failed", "error", err)
		return fmt.Errorf("redis GEOADD failed: %w", err)
	}

	return nil
}

// NearByImages fetches nearby image IDs with pagination
func (t *geoClient) GetNearByImages(
	ctx context.Context,
	long, lat, radius float64,
	userId string,
	cursor, limit int,
) ([]string, int, bool, error) {

	ctx, cancel := context.WithTimeout(ctx, 800*time.Millisecond)
	defer cancel()

	cacheKey := fmt.Sprintf("geo:nearby:%s", userId)
	sourceKey := "geo:images"

	exists, err := t.rdb.Exists(ctx, cacheKey).Result()
	if err != nil {
		t.logger.Error(
			"cacheKey not exists",
			"cacheKey", cacheKey,
			"userId", "userId",
		)
		return nil, cursor, false, err
	}

	if exists == 0 {
		_, err := t.rdb.GeoSearchStore(
			ctx,
			sourceKey, // Key--> Data selected from here
			cacheKey,  // store ---> store here for some period of time
			&redis.GeoSearchStoreQuery{
				GeoSearchQuery: redis.GeoSearchQuery{
					Longitude:  long,
					Latitude:   lat,
					Radius:     radius,
					RadiusUnit: "km",
					Sort:       "ASC",
					Count:      limit * 5,
				},
				StoreDist: true,
			},
		).Result()

		if err != nil {
			return nil, cursor, false, err
		}

		_ = t.rdb.Expire(ctx, cacheKey, 30*time.Second)
	}

	start := int64(cursor)
	end := start + int64(limit) - 1

	ids, err := t.rdb.ZRange(ctx, cacheKey, start, end).Result()
	if err != nil {
		return nil, cursor, false, err
	}

	nextCursor := cursor + len(ids)
	hasMore := len(ids) == limit

	return ids, nextCursor, hasMore, nil
}
