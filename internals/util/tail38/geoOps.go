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

// // NearByImages fetches nearby image IDs with pagination
// func (t *geoClient) NearByImages(
// 	ctx context.Context,
// 	collection string,
// 	long float32,
// 	lat float32,
// 	rangeUnit int,
// 	limit int,
// 	cursor int,
// ) ([]string, int, bool, error) {

// 	t.logger.Info(
// 		"Tile38 NEARBY with pagination",
// 		"collection", collection,
// 		"lon", long,
// 		"lat", lat,
// 		"rangeMeters", rangeUnit,
// 		"limit", limit,
// 		"cursor", cursor,
// 	)

// 	meters := rangeUnit * 1000

// 	args := []interface{}{
// 		"NEARBY",
// 		collection,
// 		"LIMIT",
// 		limit,
// 		"CURSOR",
// 		cursor,
// 		"POINT",
// 		long,
// 		lat,
// 		meters,
// 	}

// 	cmd := t.rdb.Do(ctx, args...)
// 	if err := cmd.Err(); err != nil {
// 		t.logger.Error("Tile38 NEARBY failed", "error", err.Error())
// 		return nil, 0, false, err
// 	}

// 	raw, err := cmd.Result()
// 	if err != nil {
// 		return nil, 0, false, err
// 	}

// 	t.logger.Info("Row result from tile38", "rawResult", raw)

// 	resp, ok := raw.(map[string]interface{})
// 	if !ok {
// 		return nil, 0, false, fmt.Errorf("unexpected tile38 response format")
// 	}

// 	// Parse IDs
// 	objects := resp["objects"].([]interface{})
// 	ids := make([]string, 0, len(objects))
// 	for _, obj := range objects {
// 		o := obj.(map[string]interface{})
// 		ids = append(ids, o["id"].(string))
// 	}

// 	// Cursor handling
// 	nextCursor := int(resp["cursor"].(float64))
// 	hasMore := nextCursor != 0

// 	return ids, nextCursor, hasMore, nil
// }
