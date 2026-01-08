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
}
