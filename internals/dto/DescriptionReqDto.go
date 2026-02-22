package dto

import "truthly/internals/model"

type DescriptionReqDto struct {
	ImageId string `json:"imageId"`
	UserId  string `json:"userId"`

	Description *string `json:"description"`
	Country     *string `json:"country"`
	State       *string `json:"state"`
	City        *string `json:"city"`
}

func ToDescriptionModel(d *DescriptionReqDto) *model.Description {
	return &model.Description{
		ImageId:     d.ImageId,
		UserId:      d.UserId,
		Description: d.Description,
		Country:     d.Country,
		State:       d.State,
		City:        d.City,
	}
}
