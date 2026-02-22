package dto

import "truthly/internals/model"

type CommentReqDto struct {
	UserId        string `json:"userId"`
	ImageId       string `json:"imageId"`
	DescriptionId string `json:"descriptionId"`
	AnalyticId    string `json:"analyticid"`

	Comment *string `json:"comment"`
}

func ToCommentModel(c *CommentReqDto) *model.Commemts {
	return &model.Commemts{
		UserId:        c.UserId,
		ImageId:       c.ImageId,
		DescriptionId: c.DescriptionId,
		AnalyticId:    c.AnalyticId,
		Comment:       c.Comment,
	}
}
