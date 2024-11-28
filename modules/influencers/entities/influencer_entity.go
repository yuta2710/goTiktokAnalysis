package entities

import (
	"github.com/yuta_2710/go-clean-arc-reviews/shared"
)

type Influencer struct {
	shared.BaseSQLModel
	Username        string `gorm:"column:username" json:"username"`
	DisplayName     string `gorm:"column:display_name" json:"displayName"`
	Bio             string `gorm:"column:bio" json:"bio"`
	Avatar          string `gorm:"column:avatar" json:"avatar`
	FollowersCount  int64  `gorm:"column:followers_count" json:"followersCount`
	FollowingsCount int64  `gorm:"column:followings_count" json:"followingsCount`
	TotalPostsCount int    `gorm:"column:total_posts_count" json:"totalPostsCount`
	TotalViewCount  int    `gorm:"column:total_view_count" json:"totalViewCount`
	Country         string `gorm:"column:country" json:"country"`
}

// Forked https://github.com/kluctl/go-embed-python
