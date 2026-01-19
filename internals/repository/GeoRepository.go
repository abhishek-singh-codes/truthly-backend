package repository

import (
	"context"
	"truthly/internals/util/tail38"
)

type GeoRepository interface {
	SaveImageLocation(
		ctx context.Context,
		imageId string,
		long float64,
		lat float64,
	) error

	NearByImages(
		ctx context.Context,
		collection string,
		long, lat float64,
		radius float64,
		userId string,
		cursor, limit int,
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
	long float64,
	lat float64,

) error {
	return r.client.SetPoint(
		ctx,
		"images",
		imageId,
		long,
		lat,
	)
}

func (r *tail38GeoRepository) NearByImages(
	ctx context.Context,
	collection string,
	long, lat float64,
	radius float64,
	userId string,
	cursor, limit int,
) ([]string, int, bool, error) {
	return r.client.GetNearByImages(
		ctx, collection,
		long, lat, radius,
		userId, cursor, limit,
	)
}
