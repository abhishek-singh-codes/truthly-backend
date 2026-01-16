package dto

import "time"

type FeedResponseDto struct {
	Items      []FeedItemDto `json:"items"`
	Pagination PaginationDto `json:"pagination"`
}

type FeedItemDto struct {
	ImageId  string `json:"imageId"`
	ImageUrl string `json:"imageUrl"`
	Caption  string `json:"caption"`
	UserName string `json:"userName"`
	UserId   string `json:"UserId"`

	Location  LocationDto  `json:"location"`
	Analytics AnalyticsDto `json:"analytics"`

	CreatedAt time.Time `json:"createdAt"`
}

type LocationDto struct {
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
}

type AnalyticsDto struct {
	Like      int  `json:"like"`
	Share     int  `json:"share"`
	Comment   int  `json:"comment"`
	IsLike    bool `json:"isLike"`
	IsComment bool `json:"isComment"`
}

type PaginationDto struct {
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

// Feed by range
type GetFeedByRangeDto struct {
	Limit      int
	Cursor     int
	UserId     string
	RangeUnit  int
	Collection string
	Long       float32
	Lat        float32
}
