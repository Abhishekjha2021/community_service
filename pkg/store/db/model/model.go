package model

import "time"

type UserDetails struct {
	ID              int64     `gorm:"primary_key;column:id;autoIncrement"`
	UserID          string    `gorm:"column:user_id"`
	UserName        string    `gorm:"column:user_name"`
	FirstName       string    `gorm:"column:first_name"`
	MiddleName      string    `gorm:"column:middle_name"`
	LastName        string    `gorm:"column:last_name"`
	ProfileImageUrl string    `gorm:"column:profile_image_url"`
	Email           string    `gorm:"column:email"`
	UserPhone       string    `gorm:"column:user_phone"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

type Post struct {
	ID            int64     `gorm:"primary_key;column:id;autoIncrement;"`
	ChannelID     string    `gorm:"column:channel_id"`
	UserID        string    `gorm:"column:user_id"`
	Content       string    `gorm:"column:content"`
	Type          string    `gorm:"column:type"`
	ParentID      int64     `gorm:"column:parent_id"`
	LikeCount     int64     `gorm:"column:like_count"`
	BookMarkCount int64     `gorm:"column:bookmark_count"`
	IsPinned      bool      `gorm:"column:is_pinned"`
	Status        string    `gorm:"column:status"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at;default:NULL"`
}

type UserActions struct {
	ID        int64     `gorm:"primary_key;column:id;autoIncrement"`
	UserID    string    `gorm:"column:user_id"`
	PostID    int64     `gorm:"column:post_id"`
	Action    string    `gorm:"column:action"`
	Value     *bool     `gorm:"column:value"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type UserStatus struct {
	ID            int64     `gorm:"primary_key;column:id;autoIncrement"`
	UserID        string    `gorm:"column:user_id"`
	WarningsCount int64     `gorm:"column:warnings_count"`
	Status        string    `gorm:"column:status"`
	BlockedUntil  time.Time `gorm:"column:blocked_until"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}

type MasterReport struct {
	ID        int64     `gorm:"primary_key;column:id;autoIncrement"`
	Title     string    `gorm:"column:title"`
	Subtitle  string    `gorm:"column:subtitle"`
	Topic     string    `gorm:"column:topic"`
	ParentID  int64     `gorm:"column:parent_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type Forum struct {
	ID        int64     `gorm:"primary_key;column:id;autoIncrement"`
	Title     string    `gorm:"column:title"`
	Subtitle  string    `gorm:"column:sub_title"`
	ImageURL  string    `gorm:"column:image_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

type Reports struct {
	ID             int64     `gorm:"primary_key;column:id;autoIncrement"`
	PostID         string    `gorm:"column:post_id"`
	ReportedBy     string    `gorm:"column:reported_by"`
	MasterReportID int64     `gorm:"column:master_report_id"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

type ForumEventLink struct {
	ID        int64  `gorm:"primary_key;column:id;autoIncrement"`
	ForumID   int64  `gorm:"column:forum_id"`
	ChannelID string `gorm:"column:channel_id"`
}
