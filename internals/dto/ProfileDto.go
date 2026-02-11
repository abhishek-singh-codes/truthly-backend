package dto

type ProfileResponseDto struct {
	Bio      string   `json:"bio"`
	UserName string   `json:"userName"`
	Images   []string `json:"images"`
	Freinds  int64    `json:"freinds"`
}
