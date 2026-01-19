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
	Longitude float64 `form:"longitude"`
	Latitude  float64 `form:"latitude"`
}
