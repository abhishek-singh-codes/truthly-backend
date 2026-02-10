package dto

type ProfileResponseDto struct {
	Bio        string   `json:"bio"`
	UserName   string   `json:"userName"`
	Images     []string `json:"images"`
	Neighbours int64    `json:"neighbours"`
}
