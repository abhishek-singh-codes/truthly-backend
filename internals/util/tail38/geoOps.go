package tail38

import (
	"context"
	"fmt"

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
func (t *geoClient) GetImageByRange(
	ctx context.Context, 
	collection string, 
	long float64, 
	lat float64,
	cursor 
)