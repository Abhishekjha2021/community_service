package api

import "github.com/gin-gonic/gin"

type CommunityController interface {
	CreatePost(ctx *gin.Context)
	GetPosts(ctx *gin.Context)
	LikePost(ctx *gin.Context)
	DeletePost(ctx *gin.Context)
	ReportPost(ctx *gin.Context)
	AllRepliesOnPost(ctx *gin.Context)
	MarkNotificationsAsRead(ctx *gin.Context)
}
