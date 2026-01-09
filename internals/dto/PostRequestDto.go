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
	Latitude  float64 `form:"latitude"`
	Longitude float64 `form:"longitude"`
}
