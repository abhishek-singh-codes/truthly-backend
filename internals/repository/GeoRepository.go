package repository

import (
	"context"
	"truthly/internals/util/tail38"
)

type GeoRepository interface {
	SaveImageLocation(
		ctx context.Context,
		imageId string,
		lat float32,
		long float32,
	) error
}

type tail38GeoRepository struct {
	client tail38.GeoClient
}

// constructor
func GetNewTail38Repository(client tail38.GeoClient) GeoRepository {
	return &tail38GeoRepository{
		client: client,
	}
}

func (r *tail38GeoRepository) SaveImageLocation(
	ctx context.Context,
	imageId string,
	lat float32,
	long float32,
) error {
	return r.client.SetPoint(
		ctx,
		"images",
		imageId,
		lat,
		long,
	)
}
