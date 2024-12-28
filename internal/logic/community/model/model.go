package model

import (
	"time"

	"github.com/Abhishekjha321/community_service/internal/common"
)

type RequestCreatePost struct {
	ChannelID   string `json:"channel_id"`
	Content     string `json:"content"`
	CommentType string `json:"comment_type"`
	ParentID    int    `json:"parent_id"`
}

type CreatePostData struct {
	UserName        string    `json:"user_name"`
	ProfileImageUrl string    `json:"profile_image_url"`
	UserID          string    `json:"user_id"`
	PostID          int       `json:"post_id"`
	ChannelID       string    `json:"channel_id"`
	Content         string    `json:"content"`
	CommentType     string    `json:"comment_type"`
	ParentID        int       `json:"parent_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
type ResponseCreatePost struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Data    *CreatePostData `json:"data"`
}

type ResponseGetEventPostEngagementData struct {
	Code    string               `json:"code"`
	Message string               `json:"message"`
	Data    EventLikeCommentData `json:"data"`
}

type EventLikeCommentData struct {
	LikeCount    int64 `json:"likeCount"`
	CommentCount int64 `json:"commentCount"`
}

type ResponsePostWithReplies struct {
	Code       string             `json:"code"`
	Message    string             `json:"message"`
	Data       []*PostWithReplies `json:"data"`
	Pagination common.Pagination  `json:"pagination"`
}

type PostWithReplies struct {
	ID        int                  `json:"id"`
	Avatar    string               `json:"avatar"`
	UserName  string               `json:"user_name"`
	Content   string               `json:"content"`
	Type      string               `json:"type"`
	LikeCount int                  `json:"like_count"`
	Status    string               `json:"status"`
	UserID    string               `json:"user_id"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	Replies   []common.ReponsePost `json:"replies"`
}

type RequestReportPost struct {
	PostID         string `json:"post_id"`
	MasterReportID int    `json:"master_report_id"`
}

type RepliesOnPost struct {
	// UserID          string    `json:"user_id"`
	PostID          string     `json:"post_id"`
	UserName        string     `json:"user_name"`
	ProfileImageUrl string     `json:"profile_image_url"`
	Content         string     `json:"content"`
	Type            string     `json:"type"`
	LikeCount       int        `json:"like_count"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Replies         []*Replies `json:"replies" gorm:"-"`
}

type Replies struct {
	PostID          string    `json:"post_id"`
	UserName        string    `json:"user_name"`
	ProfileImageUrl string    `json:"profile_image_url"`
	Content         string    `json:"content"`
	Type            string    `json:"type"`
	LikeCount       int       `json:"like_count"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
type ResponseRepliesOnPost struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Data       RepliesOnPost     `json:"data"`
	Pagination common.Pagination `json:"pagination"`
}

type GetPostByPostId struct {
	PostID          string    `json:"post_id"`
	UserName        string    `json:"user_name"`
	ProfileImageUrl string    `json:"profile_image_url"`
	Content         string    `json:"content"`
	Type            string    `json:"type"`
	LikeCount       int       `json:"like_count"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Post struct {
	PostID          string    `json:"post_id"`
	UserName        string    `json:"user_name"`
	ProfileImageUrl string    `json:"profile_image_url"`
	Content         string    `json:"content"`
	Type            string    `json:"type"`
	LikeCount       int       `json:"like_count"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserId          string    `json:"user_id"`
	IsLiked         bool      `json:"is_liked"`
}

type Result struct {
	IsUnread bool `json:"is_unread"`
}
type ResponseMarkAsRead struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    Result `json:"data"`
}

// -------------------repo structs-------------------

type UserAction struct {
	Action string
	Value  bool
}

type LikeCommentCountOnChannel struct {
	LikeCount    int64 `json:"like_count"`
	CommentCount int64 `json:"comment_count"`
}

type PostUserDetails struct {
	ID              int    `json:"id"`
	UserID          string `json:"user_id"`
	FirstName       string `json:"first_name"`
	MiddleName      string `json:"middle_name"`
	LastName        string `json:"last_name"`
	UserPhone       string `json:"user_phone"`
	ProfileImageUrl string `json:"profile_image_url"`
	Action          string `json:"action"`
	Value           bool   `json:"value"`
}

type ReplyPost struct {
	ID              int64
	UserID          string
	UserPhone       string
	FirstName       string
	MiddleName      string
	LastName        string
	ProfileImageUrl string
	Content         string
	Type            string
	Status          string
	LikeCount       int64
	CreatedAt       string
	UpdatedAt       string
}

type ReportResponse struct {
	ParentId       int64  `json:"parent_id"`
	ParentTitle    string `json:"parent_title"`
	ParentTopic    string `json:"parent_topic"`
	ParentSubTitle string `json:"parent_subtitle"`
	ChildId        int64  `json:"child_id"`
	ChildTitle     string `json:"child_title"`
	ChildTopic     string `json:"child_topic"`
	ChildSubTitle  string `json:"child_subtitle"`
}
