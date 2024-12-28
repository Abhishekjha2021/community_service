package proto

type RequestCreatePost struct {
	ChannelID   string `json:"channel_id"`
	Content     string `json:"content"`
	CommentType string `json:"comment_type"`
	ParentID    int64  `json:"parent_id"`
}

type ResponseCreatePost struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    ResponseCreatePostData `json:"data"`
}

type ResponseCreatePostData struct {
	UserName        string `json:"user_name"`
	ProfileImageURL string `json:"profile_image_url"`
	UserID          string `json:"user_id"`
	UserPhone       string `json:"user_phone"`
	PostID          int64  `json:"post_id"`
	ChannelID       string `json:"channel_id"`
	Content         string `json:"content"`
	CommentType     string `json:"comment_type"`
	ParentID        int64  `json:"parent_id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type ResponseAllRepliesOnPost struct {
	Code       string                           `json:"code"`
	Message    string                           `json:"message"`
	Data       ResponseAllRepliesOnPostData    `json:"data"`
	Pagination ResponseAllRepliesOnPostPagination `json:"pagination"`
}

type ResponseAllRepliesOnPostData struct {
	PostID           string                              `json:"post_id"`
	UserName         string                              `json:"user_name"`
	ProfileImageURL  string                              `json:"profile_image_url"`
	Content          string                              `json:"content"`
	Type             string                              `json:"type"`
	LikeCount        int64                               `json:"like_count"`
	Status           string                              `json:"status"`
	CreatedAt        string                              `json:"created_at"`
	UpdatedAt        string                              `json:"updated_at"`
	IsLiked          bool                                `json:"is_liked"`
	UserID           string                              `json:"user_id"`
	UserPhone        string                              `json:"user_phone"`
	BookmarkCount    int64                               `json:"bookmark_count"`
	IsBookmarked     bool                                `json:"is_bookmarked"`
	Replies          []ResponseAllRepliesOnPostReplies   `json:"replies"`
}

type ResponseAllRepliesOnPostReplies struct {
	PostID          string `json:"post_id"`
	UserName        string `json:"user_name"`
	ProfileImageURL string `json:"profile_image_url"`
	Content         string `json:"content"`
	Type            string `json:"type"`
	LikeCount       int64  `json:"like_count"`
	Status          string `json:"status"`
	IsLiked         bool   `json:"is_liked"`
	UserID          string `json:"user_id"`
	UserPhone       string `json:"user_phone"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

type ResponseAllRepliesOnPostPagination struct {
	CurrentPage            int64 `json:"current_page"`
	TotalPages             int64 `json:"total_pages"`
	SinglePageRecordCount  int64 `json:"single_page_record_count"`
	TotalRecordCount       int64 `json:"total_record_count"`
}

type ResponseGetPosts struct {
	Code       string                     `json:"code"`
	Message    string                     `json:"message"`
	Data       []ResponseGetPostsPostData `json:"data"`
	Pagination ResponseGetPostsPagination `json:"pagination"`
}

type ResponseGetPostsPagination struct {
	CurrentPage           int64 `json:"current_page"`
	TotalPages            int64 `json:"total_pages"`
	SinglePageRecordCount int64 `json:"single_page_record_count"`
	TotalRecordCount      int64 `json:"total_record_count"`
}

type ResponseGetPostsReply struct {
	ID         int64  `json:"id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	UserID     string `json:"user_id"`
	LikeCount  int64  `json:"like_count"`
	Status     string `json:"status"`
	Avatar     string `json:"avatar"`
	UserName   string `json:"user_name"`
	UserPhone  string `json:"user_phone"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	IsLiked    bool   `json:"is_liked"`
}

type ResponseGetPostsPostData struct {
	ID            int64                   `json:"id"`
	Avatar        string                  `json:"avatar"`
	UserName      string                  `json:"user_name"`
	UserPhone     string                  `json:"user_phone"`
	Content       string                  `json:"content"`
	Type          string                  `json:"type"`
	LikeCount     int64                   `json:"like_count"`
	Status        string                  `json:"status"`
	UserID        string                  `json:"user_id"`
	CreatedAt     string                  `json:"created_at"`
	UpdatedAt     string                  `json:"updated_at"`
	IsLiked       bool                    `json:"is_liked"`
	RepliesCount  int64                   `json:"replies_count"`
	BookmarkCount int64                   `json:"bookmark_count"`
	IsBookmarked  bool                    `json:"is_bookmarked"`
	IsPinned      bool                    `json:"is_pinned"`
	Replies       []ResponseGetPostsReply `json:"replies"`
}

type ResponseMarkNotificationsAsRead struct {
	Code    string                          `json:"code"`
	Message string                          `json:"message"`
	Data    ResponseMarkNotificationsAsReadData `json:"data"`
}

type ResponseMarkNotificationsAsReadData struct {
	IsUnread bool `json:"is_unread"`
}
