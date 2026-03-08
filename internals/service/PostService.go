package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"truthly/internals/dto"
	"truthly/internals/model"
	"truthly/internals/repository"
)

type PostService interface {
	UploadPost(ctx context.Context, postReq *dto.PostRequestDto, userId string) (*dto.ResponseDto[any], error)
}

type postService struct {
	logger          *slog.Logger
	analyticsRepo   repository.AnalyticRepository
	commentRepo     repository.CommentRepository
	descriptionRepo repository.DescriptionRepository
	imageRepo       repository.ImageRepository
	s3Uploader      *S3Uploader
	geoRepository   repository.GeoRepository
}

func GetPostService(
	logger *slog.Logger,
	analyticsRepo repository.AnalyticRepository,
	commentRepo repository.CommentRepository,
	descriptionRepo repository.DescriptionRepository,
	imageRepo repository.ImageRepository,
	s3Uploader *S3Uploader,
	geoRepository repository.GeoRepository,
) PostService {
	return &postService{
		logger:          logger,
		analyticsRepo:   analyticsRepo,
		commentRepo:     commentRepo,
		descriptionRepo: descriptionRepo,
		imageRepo:       imageRepo,
		s3Uploader:      s3Uploader,
		geoRepository:   geoRepository,
	}
}

func (s *postService) UploadPost(ctx context.Context, postReq *dto.PostRequestDto, userId string) (*dto.ResponseDto[any], error) {
	s.logger.Info("Uploading image...", "userId", userId)

	// 1. validate file
	if postReq.FileHeader == nil {
		s.logger.Error("no image file provided in request")

		return &dto.ResponseDto[any]{
			Status:    "failed",
			Message:   "image file is required",
			ResultObj: nil,
		}, fmt.Errorf("image file is required")
	}

	fileHeader := postReq.FileHeader

	// unique file name
	fileName := fmt.Sprintf("uploads/%d-%s", time.Now().Unix(), fileHeader.Filename)

	//1. Upload the image on s3 bucket and get url
	imgUrl, err := s.s3Uploader.UploadImage(fileHeader, fileName)
	if err != nil {
		s.logger.Error(err.Error())
		return &dto.ResponseDto[any]{
			Status:    "failed",
			Message:   err.Error(),
			ResultObj: nil,
		}, err
	}

	//2. Insert row in Image table
	img := &model.Image{
		ImageUrl:  imgUrl,
		UserId:    userId,
		Latitude:  postReq.Latitude,
		Longitude: postReq.Longitude,
	}

	imgRes, err := s.imageRepo.InsertNewImage(ctx, img)
	if err != nil {
		s.logger.Error(err.Error())
		return &dto.ResponseDto[any]{
			Status:    "failed",
			Message:   err.Error(),
			ResultObj: nil,
		}, err
	}

	// 3. Insert the description row in table like city, state etc

	description := &model.Description{
		ImageId: imgRes.ImageId,

		Description: postReq.Description,
		Country:     postReq.Country,
		City:        postReq.City,
		State:       postReq.State,
	}

	descRes, err := s.descriptionRepo.InsertDescription(ctx, description)
	if err != nil {
		s.logger.Error(err.Error())
		return &dto.ResponseDto[any]{
			Status:    "failed",
			Message:   err.Error(),
			ResultObj: nil,
		}, err
	}

	// 4. Analytic like, share, comment initially set 0

	analytic := &model.Analytic{
		ImageId: descRes.ImageId,
		// other field will be zero
	}

	analyticRes, err := s.analyticsRepo.InsertAnalytics(ctx, analytic)
	if err != nil {
		s.logger.Error(err.Error())
		return &dto.ResponseDto[any]{
			Status:    "failed",
			Message:   err.Error(),
			ResultObj: nil,
		}, err
	}

	// 5. Comment initally nothing
	comment := &model.Comment{
		ImageId: analyticRes.ImageId,
	}

	_, err = s.commentRepo.InsertComment(ctx, comment)
	if err != nil {
		s.logger.Error(err.Error())
		return &dto.ResponseDto[any]{
			Status:    "failed",
			Message:   err.Error(),
			ResultObj: nil,
			Error:     err.Error(),
		}, err
	}

	s.logger.Info("All data (Tables) related to this image are updated succesfully",
		"imageId", imgRes.ImageId, "url", imgRes.ImageUrl,
	)

	// call the geo location repo for adding data into tail38
	if postReq.Latitude != nil && postReq.Longitude != nil {

		err := s.geoRepository.SaveImageLocation(
			ctx,
			imgRes.ImageId,
			*postReq.Longitude,
			*postReq.Latitude,
		)

		// this is blocking code
		if err != nil {
			s.logger.Error(err.Error())
			return &dto.ResponseDto[any]{
				Status:    "failed to set image in tail38",
				Message:   err.Error(),
				ResultObj: nil,
				Error:     err.Error(),
			}, err
		}
	}

	return &dto.ResponseDto[any]{
		Status:  "success",
		Message: "Post created",
		ResultObj: map[string]interface{}{
			"imageUrl": imgRes.ImageUrl,
			"imageId":  imgRes.ImageId,
		},
	}, nil
}
