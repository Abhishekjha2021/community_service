package common

import "time"

const (
	USER_BASED_FLOW  = "user-based"
	IDEAS_BASED_FLOW = "ideas-based"
	Like_Action      = "like"
	Bookmark_Action  = "bookmark"
	POST_REPLY       = "reply"
	IDEAS_CHANNEL_ID = "community-ideas"
)

type UserInfo struct {
	UserID          string
	ProfileImageUrl string
	FirstName       string
	MiddleName      string
	LastName        string
	Email           string
	UserPhone       string
	UserName        string
	Version         int64
	IsNewUser       bool
}

type Post struct {
	ID            int64     `json:"id"`
	ChannelID     string    `json:"channel_id"`
	UserID        string    `json:"user_id"`
	Content       string    `json:"content"`
	IsPinned      bool      `json:"is_pinned"`
	Type          string    `json:"type"`
	ParentID      int64     `json:"parent_id"`
	LikeCount     int64     `json:"like_count"`
	BookmarkCount int64     `json:"bookmark_count"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}

type CommentWithReplies struct {
	CommentID        int64     `json:"comment_id"`
	CommentChannelID string    `json:"comment_channel_id"`
	CommentUserID    string    `json:"comment_user_id"`
	CommentContent   string    `json:"comment_content"`
	CommentType      string    `json:"comment_type"`
	CommentParentID  int64     `json:"comment_parent_id"`
	CommentLikeCount int64     `json:"comment_like_count"`
	CommentStatus    string    `json:"comment_status"`
	CommentCreatedAt time.Time `json:"comment_created_at"`
	CommentUpdatedAt time.Time `json:"comment_updated_at"`
	CommentDeletedAt time.Time `json:"comment_deleted_at"`
	ReplyID          int64     `json:"reply_id"`
	ReplyChannelID   string    `json:"reply_channel_id"`
	ReplyUserID      string    `json:"reply_user_id"`
	ReplyContent     string    `json:"reply_content"`
	ReplyType        string    `json:"reply_type"`
	ReplyParentID    int64     `json:"reply_parent_id"`
	ReplyLikeCount   int64     `json:"reply_like_count"`
	ReplyStatus      string    `json:"reply_status"`
	ReplyCreatedAt   time.Time `json:"reply_created_at"`
	ReplyUpdatedAt   time.Time `json:"reply_updated_at"`
	ReplyDeletedAt   time.Time `json:"reply_deleted_at"`
}

type ReponsePost struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Type      string    `json:"type"`
	UserID    string    `json:"user_id"`
	LikeCount int64     `json:"like_count"`
	Status    string    `json:"status"`
	Avatar    string    `json:"avatar"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Pagination struct {
	CurrentPage           int64 `json:"current_page"`
	TotalPages            int64 `json:"total_pages"`
	SinglePageRecordCount int64 `json:"single_page_record_count"`
	TotalRecordCount      int64 `json:"total_record_count"`
}

type AllRepliesPost struct {
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
	UserPhone       string    `json:"user_phone"`
	BookmarkCount   int64     `json:"bookmarkCount"`
	IsBookmarked    bool      `json:"isBookmarked"`
}
