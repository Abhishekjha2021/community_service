package model

import (
	"context"

	"github.com/Abhishekjha321/community_service/dto"
	"github.com/Abhishekjha321/community_service/exceptions"
	"github.com/Abhishekjha321/community_service/internal/common"
	dbModel "github.com/Abhishekjha321/community_service/pkg/store/db/model"
)

type Service interface {
	GetPosts(ctx context.Context, channelID string, userID string, limit int, currentPage int, sortBy string, bookMarksOnly bool) (*dto.ResponseGetPosts, *exceptions.Exception)
	LikePost(ctx context.Context, postID string, action string, userID string, channelID string) *exceptions.Exception
	DeletePost(ctx context.Context, postID string, userID string) *exceptions.Exception
	CreatePost(ctx context.Context, requestBody *dto.RequestCreatePost, userId string) (*dto.ResponseCreatePost, *exceptions.Exception)
	ReportPost(ctx context.Context, requestBody *RequestReportPost, userId string) *exceptions.Exception
	AllRepliesOnPost(ctx context.Context, postId string, userId string, channelID string, limit int, currentPage int, sortBy string) (*dto.ResponseAllRepliesOnPost, *exceptions.Exception)
	GetPostsCount(ctx context.Context, channelID string) (int64, *exceptions.Exception)
	GetRepliesCount(ctx context.Context, postID string) (int64, *exceptions.Exception)
	MarkAsRead(ctx context.Context, userID string, channelID string) (*dto.ResponseMarkNotificationsAsRead, *exceptions.Exception)
	HasUserReadPost(ctx context.Context, userID string, channelID string) (bool, error)
	GetUserPostsCount(ctx context.Context, channelID string, userID string, sortBy string) (int, *exceptions.Exception)
	GetUserEventPosts(ctx context.Context, channelID string, limit int, currentPage int, userID string, offset int, sortBy string) ([]common.Post, *exceptions.Exception)
}

type Repo interface {
	GetUserDetailsForPostID(ctx context.Context, postIDs []int) ([]PostUserDetails, error)
	GetEventPosts(ctx context.Context, channelID string, limit int, currentPage int, userId string, offset int, sortBy string) ([]common.Post, error)
	GetBookMarkedPosts(ctx context.Context, channelID string, limit int, currentPage int, userId string, offset int, sortBy string) ([]common.Post, error)
	ActionSpecificLikePost(ctx context.Context, postID string, action string, userID string) (string, *exceptions.Exception)
	CheckPostIDValidity(ctx context.Context, postID string, channelId string) (common.Post, error)
	DeleteSpecificPost(ctx context.Context, postID string, userID string) (string, error)
	InsertPostData(ctx context.Context, postData dbModel.Post) (*dbModel.Post, error)
	GetUserDetailsByUserId(ctx context.Context, userId string) (*dbModel.UserDetails, error)
	ReportPostData(ctx context.Context, report dbModel.Reports) (*dbModel.Reports, error)
	PopulateUserInfoTable(ctx context.Context, userInfo common.UserInfo) error
	GetAllRepliesOnPost(ctx context.Context, postId string, channelID string, limit int, currentPage int, sortBy string) ([]ReplyPost, error)
	GetEventPostsCount(ctx context.Context, channelID string) (int64, error)
	GetCommentSpecificReplyCount(ctx context.Context, postID string) (int64, error)
	GetPostByPostId(ctx context.Context, postId string) (*common.AllRepliesPost, error)
	FetchUserPostSpecificActionValue(ctx context.Context, postID int64, userID string) (bool, bool, error)
	GetUserIDByPostID(ctx context.Context, ParentId int64) (string, error)
	GetUserSpecificPostsCount(ctx context.Context, channelId string, userId string) (int, error)
	GetUserSpecificEventPosts(ctx context.Context, channelID string, limit int, currentPage int, userId string, offset int) ([]common.Post, error)
	GetRelevantReplies(ctx context.Context, commentIds []int64, sortBy string, bookMarksOnly bool) ([]common.Post, error)
	UpdateRequiredActionInPostsTable(ctx context.Context, postID string, updateExpr string, userID string, action string) *exceptions.Exception
}

type Consumer interface {
	UserInfoConsumer(msg []byte) error
}
