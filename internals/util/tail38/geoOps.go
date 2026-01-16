package tail38

import (
	"context"
	"fmt"
)

// set point functionality
func (t *geoClient) SetPoint(
	ctx context.Context,
	collection string,
	id string,
	lat float32,
	long float32,
) error {

	t.logger.Info(
		"Inserting post data into tail38",
		"imageId", id,
		"lat", lat,
		"long", long,
	)

	cmd := t.rdb.Do(
		ctx,
		"SET",
		collection,
		id,
		"POINT",
		long,
		lat,
	)

	if err := cmd.Err(); err != nil {
		t.logger.Error("tail38 SET failed", "error", err.Error())
		return fmt.Errorf("tail38 SET failed: %s", err.Error())
	}

	return nil
}

// NearByImages fetches nearby image IDs with pagination
func (t *geoClient) NearByImages(
	ctx context.Context,
	collection string,
	long float32,
	lat float32,
	rangeUnit int,
	limit int,
	cursor int,
) ([]string, int, bool, error) {

	t.logger.Info(
		"Tile38 NEARBY with pagination",
		"collection", collection,
		"lon", long,
		"lat", lat,
		"rangeMeters", rangeUnit,
		"limit", limit,
		"cursor", cursor,
	)

	meters := rangeUnit * 1000

	args := []interface{}{
		"NEARBY",
		collection,
		"LIMIT",
		limit,
		"CURSOR",
		cursor,
		"POINT",
		long,
		lat,
		meters,
	}

	cmd := t.rdb.Do(ctx, args...)
	if err := cmd.Err(); err != nil {
		t.logger.Error("Tile38 NEARBY failed", "error", err.Error())
		return nil, 0, false, err
	}

	raw, err := cmd.Result()
	if err != nil {
		return nil, 0, false, err
	}

	t.logger.Info("Row result from tile38", "rawResult", raw)

	resp, ok := raw.(map[string]interface{})
	if !ok {
		return nil, 0, false, fmt.Errorf("unexpected tile38 response format")
	}

	// Parse IDs
	objects := resp["objects"].([]interface{})
	ids := make([]string, 0, len(objects))
	for _, obj := range objects {
		o := obj.(map[string]interface{})
		ids = append(ids, o["id"].(string))
	}

	// Cursor handling
	nextCursor := int(resp["cursor"].(float64))
	hasMore := nextCursor != 0

	return ids, nextCursor, hasMore, nil
}
