package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	PublicApiV1PathPrefix      = "/community/v1"
	PrivateApiV1PathPrefix     = "/private/community/v1"
)

func AddPublicRoutes(router *gin.Engine, communityController CommunityController) {
	v1Public := router.Group(PublicApiV1PathPrefix)
	{
		v1Public.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		v1Public.GET("/post", communityController.GetPosts)
		v1Public.GET("public/post", communityController.GetPosts)
		v1Public.POST("/like", communityController.LikePost)
		v1Public.DELETE("/post", communityController.DeletePost)
		v1Public.POST("/post", communityController.CreatePost)
		v1Public.POST("/report", communityController.ReportPost)
		v1Public.GET("/replies", communityController.AllRepliesOnPost)
		v1Public.GET("public/replies", communityController.AllRepliesOnPost)
		v1Public.GET("/read_status", communityController.MarkNotificationsAsRead)
	}
}

func AddPrivateRoutes(router *gin.Engine, communityController CommunityController) {
	// add internal apis here
	v1Private := router.Group(PrivateApiV1PathPrefix)
	{
		v1Private.GET("/health", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
	}
}
