package dao

import (
	"gorm.io/gorm"
	"time"
)

type OauthAccessToken struct {
	Id                  int64 `gorm:"primaryKey"`
	UserId              int64
	AccessTokenUUID     string
	RefreshTokenUUID    string
	AccessTokenExpDate  time.Time
	RefreshTokenExpData time.Time
	Revoked             bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}
