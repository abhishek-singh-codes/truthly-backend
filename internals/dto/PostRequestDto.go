package dto

import "mime/multipart"

type PostRequestDto struct {
	// Description
	Description string `form:"description"`
	Country     string `form:"country"`
	State       string `form:"state"`
	City        string `form:"city"`

	// Image
	FileHeader *multipart.FileHeader `form:"fileHeader"`

	// lat long
	Latitude  float32 `form:"latitude"`
	Longitude float32 `form:"longitude"`
}
