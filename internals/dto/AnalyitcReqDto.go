package dto

import (
	"truthly/internals/model"
)

type AnalyticReqDto struct {
	ImageId string `json:"imageId"`

	Like    *int `json:"like"`
	Share   *int `json:"share"`
	Comment *int `json:"comment"`
}

func ToAnalyticModel(a *AnalyticReqDto) *model.Analytic {
	return &model.Analytic{
		ImageId:      a.ImageId,
		LikeCount:    *a.Like,
		ShareCount:   *a.Share,
		CommentCount: *a.Comment,
	}
}
