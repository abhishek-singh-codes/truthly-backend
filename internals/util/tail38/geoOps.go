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
		lat,
		long,
	)

	if err := cmd.Err(); err != nil {
		t.logger.Error("tail38 SET failed", "error", err.Error())
		return fmt.Errorf("tail38 SET failed: %s", err.Error())
	}

	return nil
}
