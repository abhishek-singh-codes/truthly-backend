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

	FindImagesByRange(
		ctx context.Context,
		collection string,
		long float32,
		lat float32,
		rangeUnit int,
		limit int,
		cursor int,
	) ([]string, int, bool, error)
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
		long,
		lat,
	)
}

func (r *tail38GeoRepository) FindImagesByRange(
	ctx context.Context,
	collection string,
	long float32,
	lat float32,
	rangeUnit int,
	limit int,
	cursor int,
) ([]string, int, bool, error) {

	return r.client.NearByImages(
		ctx,
		collection,
		long,
		lat,
		rangeUnit,
		limit,
		cursor,
	)
}
