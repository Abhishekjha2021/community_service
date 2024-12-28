package api

import "github.com/gin-gonic/gin"

type Controller interface {
	GetPosts(ctx *gin.Context)
	LikePost(ctx *gin.Context)
}
