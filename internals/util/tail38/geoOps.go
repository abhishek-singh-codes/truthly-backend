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

	t.logger.Info(
		"Inserting image location into Redis GEO",
		"imageId", id,
		"lat", lat,
		"long", long,
	)

	err := t.rdb.GeoAdd(ctx, collection, &redis.GeoLocation{
		Name:      id,
		Longitude: long,
		Latitude:  lat,
	}).Err()

	if err != nil {
		t.logger.Error("Redis GEOADD failed", "error", err)
		return fmt.Errorf("redis GEOADD failed: %w", err)
	}

	return nil
}

// NearByImages fetches nearby image IDs with pagination
func (t *geoClient) GetNearByImages(
	ctx context.Context,
	collection string,
	long, lat float64,
	radius float64,
	userId string,
	cursor, limit int,
) ([]string, int, bool, error) {

	ctx, cancel := context.WithTimeout(ctx, 800*time.Millisecond)
	defer cancel()

	key := fmt.Sprintf("nearby:%s", userId)

	exists, err := t.rdb.Exists(ctx, key).Result()
	if err != nil {
		return nil, cursor, false, err
	}

	if exists == 0 {
		_, err := t.rdb.GeoSearchStore(
			ctx,
			key,
			collection,
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

		_ = t.rdb.Expire(ctx, key, 30*time.Second)
	}

	start := int64(cursor)
	end := int64(cursor + limit)

	imageIds, err := t.rdb.ZRange(ctx, key, start, end).Result()
	if err != nil {
		return nil, cursor, false, err
	}

	hasMore := false
	if len(imageIds) > limit {
		hasMore = true
		imageIds = imageIds[:limit]
	}

	nextCursor := cursor + len(imageIds)

	_ = t.rdb.Expire(ctx, key, 30*time.Second)

	return imageIds, nextCursor, hasMore, nil
}
