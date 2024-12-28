package api

import (
	"strconv"
	"strings"

	"github.com/Abhishekjha321/community_service/exceptions"
	"github.com/Abhishekjha321/community_service/internal/common"
	"github.com/Abhishekjha321/community_service/internal/logic/community/model"
	logger "github.com/Abhishekjha321/community_service/log"
	proto "github.com/Abhishekjha321/community_service/proto" 
	"github.com/gin-gonic/gin"
)

const (
	X_USER_ID                       = "x-user-id"
	filterBy                        = "filter_by"
	QueryParamPage           string = "page"
	QueryParamLimit          string = "limit"
	DefaultPage                     = "1"
	DefaultPageLimit                = "10"
	QueryParamPageSize              = "pageSize"
	QueryParamSortBy                = "sort_by"
	QueryParamsBookmarksOnly        = "bookmarks_only"
	CHANNEL_ID                      = "channel_id"
	POST_ID                         = "post_id"
	ACTION                          = "action"
	LIKE                            = "like"
	UNLIKE                          = "unlike"
	BOOKMARK                        = "bookmark"
	LIMIT                           = "limit"
	CURRENT_PAGE                    = "current_page"
	SORT_BY_POSTS                   = common.USER_BASED_FLOW
)

type communityController struct {
	communityService model.Service
}

type CommunityController interface {
	CreatePost(ctx *gin.Context)
	GetPosts(ctx *gin.Context)
	LikePost(ctx *gin.Context)
	DeletePost(ctx *gin.Context)
	ReportPost(ctx *gin.Context)
	AllRepliesOnPost(ctx *gin.Context)
	MarkNotificationsAsRead(ctx *gin.Context)
}

func NewCommunityController(communityService model.Service) CommunityController {
	return &communityController{
		communityService: communityService,
	}
}

func (c *communityController) CreatePost(ctx *gin.Context) {

	var (
		requestCreatePost proto.RequestCreatePost
	)

	log := logger.GetLogInstance(ctx, "CreatePost")

	userId := ctx.GetHeader(X_USER_ID)
	if userId == "" {
		log.Errorf("[CreatePostController] x-user-id not found in header")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(exceptions.UserIDMissingErrorCode))
		return
	}

	if err := ctx.BindJSON(&requestCreatePost); err != nil {

		log.Errorf("[CreatePostController] Error occurred while binding request: %v", err)

		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.BadRequestErrorCode, "[CreatePostController] Error occured while binding json"))
		return
	}

	if requestCreatePost.ChannelID == "" {
		log.Errorf("[CreatePostController] ChannelID isn't found in request body")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.BadRequestErrorCode, "[CreatePostController] ChannelID isn't found"))
		return
	}

	response, exp := c.communityService.CreatePost(ctx.Request.Context(), &requestCreatePost, userId)

	if exp != nil {

		log.Errorf("[CreatePostController] Error ocured while creating post")
		SendApiResponseV1(ctx, nil, exp)
		return
	}

	SendApiResponseV1(ctx, response, nil)
}

func (c *communityController) GetPosts(ctx *gin.Context) {
	log := logger.GetLogInstance(ctx, "Get Posts Controller")
	userID := ctx.GetHeader(X_USER_ID)
	channelID := ctx.Query(CHANNEL_ID)
	if len(channelID) == 0 {
		log.Errorf("[GetPostsController] channel id not being sent in query params")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(exceptions.QueryParamsIncorrectErrorCode))
		return
	}
	sortBy := ctx.Query(QueryParamSortBy)
	if len(sortBy) == 0 {
		sortBy = common.USER_BASED_FLOW
	}
	bookMarksOnlyParam := ctx.Query(QueryParamsBookmarksOnly)
	bookMarksOnly := false
	if len(bookMarksOnlyParam) > 0 {
		var err error
		bookMarksOnly, err = strconv.ParseBool(strings.ToLower(bookMarksOnlyParam))
		if err != nil {
			log.Errorf("[YourFunction] Invalid boolean value for bookMarksOnly: %v", err)
			bookMarksOnly = false
		}
	}

	limit := ctx.Query(LIMIT)
	if limit == "" {
		limit = DefaultPageLimit
	}
	convertedLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Errorf("[GetPostsController] couldn't convert limit: %s to integer in query params", limit)
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.BadRequestErrorCode, "Limit is missing"))
		return
	}
	currentPage := ctx.Query(CURRENT_PAGE)
	if currentPage == "" {
		currentPage = DefaultPage
	}
	convertedCurrentPage, err := strconv.Atoi(currentPage)
	if err != nil {
		log.Errorf("[GetPostsController] couldn't convert currentPage: %s to integer in query params", currentPage)
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.BadRequestErrorCode, "CurrentPage is missing"))
		return
	}
	log.Infof("channelID: %s, limit: %s, currentPage: %s", channelID, limit, currentPage)
	res, errGetPosts := c.communityService.GetPosts(ctx, channelID, userID, convertedLimit, convertedCurrentPage, sortBy, bookMarksOnly)
	if errGetPosts != nil {
		SendApiResponseV1(ctx, nil, errGetPosts)
		return
	}

	SendApiResponseV1(ctx, res, nil)
}

func getValidAction(action string) bool {
	switch strings.ToLower(action) {
	case LIKE, UNLIKE, BOOKMARK:
		return true
	}
	return false
}

func (c *communityController) LikePost(ctx *gin.Context) {

	log := logger.GetLogInstance(ctx, "Like Post ")

	userID := ctx.GetHeader(X_USER_ID)
	if userID == "" {
		log.Errorf("[LikePostController] x-user-id not found in header ")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(exceptions.UserIDMissingErrorCode))
		return
	}
	postID := ctx.Query(POST_ID)
	if postID == "" {
		log.Errorf("[LikePostController] Post Id not found in query params")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.NoDataFoundErrorCode, "Post Id not found in query params"))
		return
	}

	action := ctx.Query(ACTION)
	if !getValidAction(action) {
		log.Errorf("[LikePostController] Action sent in query params is not correct")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.QueryFailedErrorCode, "Action sent in query params, is not correct"))
		return
	}

	channelID := ctx.Query(CHANNEL_ID)

	// Information of input parameters on which we proceed with the service layer and db calls
	log.Infof("[LikePostController] postID: %s, action: %s, userID: %s", postID, action, userID)

	err := c.communityService.LikePost(ctx, postID, strings.ToLower(action), userID, channelID)
	if err != nil {
		SendApiResponseV1(ctx, nil, err)
		return
	}
	SendApiResponseV1(ctx, &SuccessResp{Code: "00000", Message: "Success"}, err)
}

func (c *communityController) DeletePost(ctx *gin.Context) {
	log := logger.GetLogInstance(ctx, "Delete Post Controller")
	userID := ctx.GetHeader(X_USER_ID)
	if userID == "" {
		log.Errorf("[DeletePostController] x-user-id not found in header ")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(exceptions.UserIDMissingErrorCode))
		return
	}
	postID := ctx.Query(POST_ID)
	if len(postID) == 0 {
		log.Errorf("[DeletePostController] Post id not being sent in query params")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(
			exceptions.PostIdErrorCode))
		return
	}
	log.Infof("[DeletePostController] postID: %s", postID)
	err := c.communityService.DeletePost(ctx, postID, userID)
	if err != nil {
		SendApiResponseV1(ctx, nil, err)
		return
	}
	SendApiResponseV1(ctx, &SuccessResp{
		Code:    "00000",
		Message: "Success",
	}, nil)
}

func (c *communityController) ReportPost(ctx *gin.Context) {
	var (
		requestReportPost model.RequestReportPost
	)

	log := logger.GetLogInstance(ctx, "ReportPost")

	userId := ctx.GetHeader(X_USER_ID)
	if userId == "" {
		log.Errorf("[ReportPostController] x-user-id not found in header for report post")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(exceptions.UserIDMissingErrorCode))
		return
	}

	if err := ctx.BindJSON(&requestReportPost); err != nil {

		log.Errorf("[ReportPostController] Error occurred while binding request for report post: %v", err)

		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.BadRequestErrorCode, "Error occured while binding json"))
		return
	}

	exp := c.communityService.ReportPost(ctx.Request.Context(), &requestReportPost, userId)

	if exp != nil {

		log.Errorf("[ReportPostController] Error ocured while reporting a post")
		SendApiResponseV1(ctx, nil, exp)
		return
	}

	SendApiResponseV1(ctx, &SuccessResp{Code: "00000", Message: "Success"}, nil)

}

func (c *communityController) AllRepliesOnPost(ctx *gin.Context) {

	log := logger.GetLogInstance(ctx, "AllRepliesOnPost")

	userId := ctx.GetHeader(X_USER_ID)
	postId := ctx.Query(POST_ID)
	channelID := ctx.Query(CHANNEL_ID)
	limit := ctx.Query(LIMIT)
	currentPage := ctx.Query(CURRENT_PAGE)

	if currentPage == "" {
		currentPage = DefaultPage
	}
	convertedCurrentPage, err := strconv.Atoi(currentPage)
	if err != nil {
		log.Errorf("[AllRepliesOnPostController] couldn't convert currentPage: %s to integer in query params", currentPage)
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.BadRequestErrorCode, "CurrentPage is missing"))
		return
	}

	if limit == "" {
		limit = DefaultPageLimit
	}
	convertedLimit, err := strconv.Atoi(limit)
	if err != nil {
		log.Errorf("[AllRepliesOnPost] couldn't convert limit: %s to integer in query params", limit)
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.BadRequestErrorCode, "Limit is missing"))
		return
	}

	if postId == "" {
		log.Errorf("[AllRepliesOnPost] Post Id not found in query params for replis on post")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(
			exceptions.PostIdErrorCode))
		return
	}

	if channelID == "" {
		log.Errorf("[AllRepliesOnPost] Channel Id not found in query params for replis on post")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(exceptions.QueryParamsIncorrectErrorCode))
		return
	}

	sortBy := ctx.Query(QueryParamSortBy)
	if len(sortBy) == 0 {
		sortBy = common.USER_BASED_FLOW
	}

	res, exp := c.communityService.AllRepliesOnPost(ctx, postId, userId, channelID, convertedLimit, convertedCurrentPage, sortBy)
	if exp != nil {
		log.Errorf("[AllRepliesOnPost] Error ocured while getting replies on a post")
		SendApiResponseV1(ctx, nil, exp)
		return
	}

	SendApiResponseV1(ctx, res, nil)

}

func (c *communityController) MarkNotificationsAsRead(ctx *gin.Context) {
	channelID := ctx.Query("channel_id")
	userID := ctx.GetHeader(X_USER_ID)
	log := logger.GetLogInstance(ctx, "MarkNotificationsAsRead")

	if userID == "" {
		log.Errorf("[MarkNotificationsAsReadController] x-user-id not found in header")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCode(exceptions.UserIDMissingErrorCode))
		return
	}

	if channelID == "" {
		log.Errorf("[MarkNotificationsAsReadController] Missing channelID in query parameters")
		SendApiResponseV1(ctx, nil, exceptions.GetExceptionByErrorCodeWithCustomMessage(
			exceptions.BadRequestErrorCode, "Missing channelID"))
		return
	}

	res, err := c.communityService.MarkAsRead(ctx, userID, channelID)
	if err != nil {
		log.Errorf("[MarkNotificationsAsReadController] Error occurred while marking notifications as read,Error: %v", err)
		SendApiResponseV1(ctx, nil, err)
		return
	}

	SendApiResponseV1(ctx, res, nil)
}
