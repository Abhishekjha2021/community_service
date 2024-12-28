package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Abhishekjha321/community_service/exceptions"
	"github.com/Abhishekjha321/community_service/internal/common"
	"github.com/Abhishekjha321/community_service/storage/cache"
	"github.com/go-redis/redis"

	// "github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	model "github.com/Abhishekjha321/community_service/internal/logic/community/model"

	"github.com/Abhishekjha321/community_service/dto"
	logger "github.com/Abhishekjha321/community_service/log"
	dbModel "github.com/Abhishekjha321/community_service/pkg/store/db/model"
)

const (
	APISuccessCode              = "00000"
	APISuccessMessage           = "SUCCESS"
	DEFAULT_COUNTRY_MOBILE_CODE = "+63"
	like                        = common.Like_Action
	deleted                     = "DELETED"
	IDEAS_CHANNEL_ID            = common.IDEAS_CHANNEL_ID
	bookmark                    = common.Bookmark_Action
	reply                       = common.POST_REPLY
)

type service struct {
	repo        model.Repo
	redisClient cache.CacheBase
	// clients     *client.ClientImpl
}

func dereferencePostData(posts []*dto.ResponseGetPostsPostData) []dto.ResponseGetPostsPostData {
	var result []dto.ResponseGetPostsPostData
	for _, post := range posts {
		result = append(result, *post)
	}
	return result
}

func NewService(repo model.Repo, redisClient cache.CacheBase) model.Service {

	return &service{
		repo:        repo,
		redisClient: redisClient,
		// clients:     clients,
	}
}

func (s *service) CreatePost(ctx context.Context, requestBody *dto.RequestCreatePost, userId string) (*dto.ResponseCreatePost, *exceptions.Exception) {
	log := logger.GetLogInstance(ctx, "CreatePost")

	post := dbModel.Post{
		UserID:    userId,
		ChannelID: requestBody.ChannelID,
		Content:   requestBody.Content,
		Type:      strings.ToUpper(requestBody.CommentType),
		ParentID:  requestBody.ParentID,
		Status:    "PUBLISHED",
	}

	reply, err := s.repo.InsertPostData(ctx, post)
	if err != nil {
		log.Errorf("[CreatePost] Error while creating post, err: %s", err)
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}

	userDetails, err := s.repo.GetUserDetailsByUserId(ctx, userId)
	if err != nil {
		log.Errorf("[CreatePost] Error while fetching user details")
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}

	//creating redis key
	if requestBody.ParentID != 0 { // This means it's a reply

		// Fetch the UserID of the parent post
		parentUserID, err := s.repo.GetUserIDByPostID(ctx, requestBody.ParentID)
		if err != nil {
			log.Errorf("[CreatePost] Error while fetching parent user ID, err: %s", err)
			return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
		}

		redisKey := fmt.Sprintf("community_comment_unread:%s:%s", parentUserID, requestBody.ChannelID)
		expiration := 7 * 24 * time.Hour
		KeyExpiryerr := s.redisClient.SetExpiringKey(ctx, redisKey, "true", expiration)
		if KeyExpiryerr != nil {
			log.Errorf("[CreatePost] Error while setting Redis key, err: %s", KeyExpiryerr)
			return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
		}
	}
	userName := userDetails.FirstName
	if userDetails.MiddleName != "" {
		userName += " " + userDetails.MiddleName
	}
	if userDetails.LastName != "" {
		userName += " " + userDetails.LastName
	}
	var result = &dto.ResponseCreatePostData{
		UserName:        userName,
		UserPhone:       userDetails.UserPhone,
		ProfileImageURL: userDetails.ProfileImageUrl,
		UserID:          userId,
		PostID:          reply.ID,
		ChannelID:       reply.ChannelID,
		Content:         reply.Content,
		CommentType:     reply.Type,
		ParentID:        reply.ParentID,
		CreatedAt:       fmt.Sprint(reply.CreatedAt.Unix()),
		UpdatedAt:       fmt.Sprint(reply.UpdatedAt.Unix()),
	}

	return &dto.ResponseCreatePost{
		Code:    APISuccessCode,
		Message: APISuccessMessage,
		Data:    *result,
	}, nil
}

func (s *service) GetPostsCount(ctx context.Context, ChannelID string) (int64, *exceptions.Exception) {
	log := logger.GetLogInstance(ctx, "GetPostsCountService")
	count, err := s.repo.GetEventPostsCount(ctx, ChannelID)
	if err != nil {
		log.Errorf("[GetPostsCountService] Couldn't get comment count for corresponding channel id: %s", ChannelID)
		return 0, exceptions.GetExceptionByErrorCode(exceptions.BadRequestErrorCode)
	}
	return count, nil
}

func (s *service) GetUserPostsCount(ctx context.Context, ChannelID string, userId string, sortBy string) (int, *exceptions.Exception) {
	log := logger.GetLogInstance(ctx, "GetUserPostsCountService")
	if strings.ToLower(sortBy) == common.IDEAS_BASED_FLOW {
		log.Infof("[GetUserPostsCountService] No need to get comment count for corresponding channel id: %s and user id: %s because it is ideas based flow", ChannelID, userId)
		return 0, nil
	}
	count, err := s.repo.GetUserSpecificPostsCount(ctx, ChannelID, userId)
	log.Infof("[GetUserPostsCountService] user specific posts count: %d for userId: %s", count, userId)
	if err != nil {
		log.Errorf("[GetUserPostsCountService] Couldn't get comment count for corresponding channel id: %s and user id: %s", ChannelID, userId)
		return count, exceptions.GetExceptionByErrorCode(exceptions.BadRequestErrorCode)
	}
	return count, nil
}

func (s *service) GetUserEventPosts(ctx context.Context, ChannelID string, limit int, currentPage int, userID string, offset int, sortBy string) ([]common.Post, *exceptions.Exception) {
	log := logger.GetLogInstance(ctx, "GetUserEventPostsService")
	if strings.ToLower(sortBy) == common.IDEAS_BASED_FLOW {
		log.Infof("[GetUserEventPosts] No need to get comment count for corresponding channel id: %s and user id: %s because it is ideas based flow", ChannelID, userID)
		return nil, nil
	}
	userPosts, err := s.repo.GetUserSpecificEventPosts(ctx, ChannelID, limit, currentPage, userID, offset)
	if err != nil {
		log.Errorf("[GetUserEventPosts] couldn't fetch user specific posts for corresponding channel id: %s and userId: %s", ChannelID, userID)
	}
	return userPosts, nil
}

func (s *service) GetPosts(ctx context.Context, ChannelID string, userID string, limit int, currentPage int, sortBy string, bookMarksOnly bool) (*dto.ResponseGetPosts, *exceptions.Exception) {
	var (
		response      *dto.ResponseGetPosts
		filteredPosts []*dto.ResponseGetPostsPostData
		log           = logger.GetLogInstance(ctx, "GetPostsService")
	)
	userPostsCount, errUserPostsCount := s.GetUserPostsCount(ctx, ChannelID, userID, sortBy)
	if errUserPostsCount != nil {
		log.Errorf("[GetPostsService] couldn't fetch user specific posts count for corresponding channel id: %s and userId: %s", ChannelID, userID)
	}
	totalPostsRequired := (currentPage - 1) * limit
	limitOtherPosts := userPostsCount - totalPostsRequired
	var userSpecificPosts []common.Post
	if limitOtherPosts > 0 {
		offsetInternal := totalPostsRequired
		userPosts, _ := s.GetUserEventPosts(ctx, ChannelID, limit, currentPage, userID, offsetInternal, sortBy)
		userSpecificPosts = userPosts
	} else {
		limitOtherPosts = 0
	}
	offsetOtherPosts := ((currentPage - 1) * limit) - userPostsCount
	if offsetOtherPosts < 0 {
		offsetOtherPosts = 0
	}
	var posts []common.Post
	if (limit - limitOtherPosts) > 0 {
		var fetchedPosts []common.Post
		var dbError error
		if bookMarksOnly {
			posts, err := s.repo.GetBookMarkedPosts(ctx, ChannelID, limit-limitOtherPosts, currentPage, userID, offsetOtherPosts, sortBy)
			fetchedPosts = posts
			dbError = err
		} else {
			posts, err := s.repo.GetEventPosts(ctx, ChannelID, limit-limitOtherPosts, currentPage, userID, offsetOtherPosts, sortBy)
			fetchedPosts = posts
			dbError = err
		}
		if dbError != nil {
			log.Errorf("[GetPostsService] not found any posts for corresponding channel id: %s", ChannelID)
			return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
		}
		if fetchedPosts == nil {
			log.Errorf("[GetPostsService] Not found any posts for corresponding channel id: %s", ChannelID)
			return response, nil
		}
		posts = append(userSpecificPosts, fetchedPosts...)
	} else {
		posts = userSpecificPosts
	}
	var commentIds []int64
	for _, post := range posts {
		commentIds = append(commentIds, post.ID)
	}
	replies, _ := s.repo.GetRelevantReplies(ctx, commentIds, sortBy, bookMarksOnly)

	var postIds []int

	for _, p := range posts {
		postIds = append(postIds, int(p.ID))
	}
	for _, p := range replies {
		postIds = append(postIds, int(p.ID))
	}

	userData, err := s.repo.GetUserDetailsForPostID(ctx, postIds)
	if err != nil {
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}

	userDataMap := make(map[string]model.PostUserDetails)

	for _, ud := range userData {
		userDataMap[ud.UserID] = ud
	}

	// Created maps to hold comments and replies
	commentMap := make(map[int64]*dto.ResponseGetPostsPostData)
	replyMap := make(map[int64][]*dto.ResponseGetPostsReply)

	// Created a list to preserve order
	var orderedComments []*dto.ResponseGetPostsPostData

	for _, post := range posts {
		if _, exists := commentMap[post.ID]; !exists {
			likeStatus, bookmarkStatus, errLikeStatus := s.repo.FetchUserPostSpecificActionValue(ctx, post.ID, userID)
			if errLikeStatus != nil {
				likeStatus = false
			}
			commentIDStr := strconv.FormatInt(post.ID, 10)
			repliesCount, errReplyCountttt := s.repo.GetCommentSpecificReplyCount(ctx, commentIDStr)
			if errReplyCountttt != nil {
				log.Errorf("[GetPostsService] Unable to fetch replies count for postId: %d", post.ID)
			}
			userName := userDataMap[post.UserID].FirstName

			if userDataMap[post.UserID].MiddleName != "" {
				userName += " " + userDataMap[post.UserID].MiddleName
			}
			if userDataMap[post.UserID].LastName != "" {
				userName += " " + userDataMap[post.UserID].LastName
			}
			comment := &dto.ResponseGetPostsPostData{
				ID:            post.ID,
				UserID:        post.UserID,
				Avatar:        userDataMap[post.UserID].ProfileImageUrl,
				UserName:      userName,
				UserPhone:     userDataMap[post.UserID].UserPhone,
				Content:       post.Content,
				Type:          post.Type,
				LikeCount:     post.LikeCount,
				Status:        post.Status,
				CreatedAt:     fmt.Sprint(post.CreatedAt.Unix()),
				UpdatedAt:     fmt.Sprint(post.UpdatedAt.Unix()),
				RepliesCount:  repliesCount,
				IsLiked:       likeStatus,
				BookmarkCount: post.BookmarkCount,
				IsBookmarked:  bookmarkStatus,
				IsPinned:      post.IsPinned,
			}
			commentMap[post.ID] = comment
			orderedComments = append(orderedComments, comment)
		}
	}

	for _, reply := range replies {
		if reply.ParentID != 0 {
			likeStatus, _, errLikeStatus := s.repo.FetchUserPostSpecificActionValue(ctx, reply.ID, userID)
			if errLikeStatus != nil {
				likeStatus = false
			}
			userName := userDataMap[reply.UserID].FirstName
			if userDataMap[reply.UserID].MiddleName != "" {
				userName += " " + userDataMap[reply.UserID].MiddleName
			}
			if userDataMap[reply.UserID].LastName != "" {
				userName += " " + userDataMap[reply.UserID].LastName
			}
			replyDetail := &dto.ResponseGetPostsReply{
				ID:        reply.ID,
				UserID:    reply.UserID,
				Avatar:    userDataMap[reply.UserID].ProfileImageUrl,
				UserName:  userName,
				UserPhone: userDataMap[reply.UserID].UserPhone,
				Content:   reply.Content,
				Type:      reply.Type,
				LikeCount: reply.LikeCount,
				Status:    reply.Status,
				CreatedAt: fmt.Sprint(reply.CreatedAt.Unix()),
				UpdatedAt: fmt.Sprint(reply.UpdatedAt.Unix()),
				IsLiked:   likeStatus,
			}

			if replies, exists := replyMap[reply.ParentID]; exists {
				replyMap[reply.ParentID] = append(replies, replyDetail)
			} else {
				replyMap[reply.ParentID] = []*dto.ResponseGetPostsReply{replyDetail}
			}
		}
	}

	for ParentID, replies := range replyMap {
		sort.Slice(replies, func(i, j int) bool {
			if replies[i].Status == "DELETED" && replies[j].Status != "DELETED" {
				return false
			}
			if replies[i].Status != "DELETED" && replies[j].Status == "DELETED" {
				return true
			}
			iUpdatedAt, _ := strconv.ParseInt(replies[i].UpdatedAt, 10, 64)
			jUpdatedAt, _ := strconv.ParseInt(replies[j].UpdatedAt, 10, 64)
			return iUpdatedAt < jUpdatedAt
		})

		if len(replies) > 3 {
			replyMap[ParentID] = replies[:3]
		}
	}

	for _, commentDetail := range orderedComments {
		postWithReplies := dto.ResponseGetPostsPostData{
			UserID:        commentDetail.UserID,
			ID:            commentDetail.ID,
			Avatar:        commentDetail.Avatar,
			UserName:      commentDetail.UserName,
			UserPhone:     commentDetail.UserPhone,
			Content:       commentDetail.Content,
			Type:          commentDetail.Type,
			LikeCount:     commentDetail.LikeCount,
			Status:        commentDetail.Status,
			CreatedAt:     commentDetail.CreatedAt,
			UpdatedAt:     commentDetail.UpdatedAt,
			IsLiked:       commentDetail.IsLiked,
			RepliesCount:  commentDetail.RepliesCount,
			BookmarkCount: commentDetail.BookmarkCount,
			IsBookmarked:  commentDetail.IsBookmarked,
			IsPinned:      commentDetail.IsPinned,
			// Replies:       replyMap[commentDetail.Id],
			Replies: dereferenceReplies(replyMap[commentDetail.ID]),
		}

		filteredPosts = append(filteredPosts, &postWithReplies)
	}

	// fetching redis key
	redisKey := fmt.Sprintf("community_comment_unread:%s:%s", userID, ChannelID)
	_, err = s.redisClient.GetKey(ctx, redisKey)
	if err != nil {
		log.Errorf("[GetPostsService] Error retrieving Redis key %s: %v", redisKey, err)
	}

	// deleting redis key
	deleteErr := s.redisClient.DeleteKey(ctx, redisKey)
	if deleteErr != nil {
		log.Errorf("[GetPostsService] Error while deleting Redis key %s: %v", redisKey, deleteErr)
	} else {
		fmt.Println("Deleted Redis key:", redisKey)
	}

	// check if the key still exists
	_, err = s.redisClient.GetKey(ctx, redisKey)
	if err != nil {
		fmt.Println("Redis key successfully deleted:", redisKey)
	} else {
		fmt.Println("Redis key still exists after deletion:", redisKey)
	}

	recordsCount, errCount := s.GetPostsCount(ctx, ChannelID)
	if errCount != nil {
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}

	pagination := NewPagination(int64(currentPage), int64(limit), recordsCount)

	return &dto.ResponseGetPosts{
		Code:       APISuccessCode,
		Message:    APISuccessMessage,
		Data:       dereferencePostData(filteredPosts),
		Pagination: *pagination,
	}, nil

}

func dereferenceReplies(replies []*dto.ResponseGetPostsReply) []dto.ResponseGetPostsReply {
	var result []dto.ResponseGetPostsReply
	for _, reply := range replies {
		result = append(result, *reply)
	}
	return result
}

func NewPagination(currentPage, pageSize, totalRecordCount int64) *dto.ResponseGetPostsPagination {
	totalPages := (totalRecordCount + pageSize - 1) / pageSize
	var singlePageRecordCount int64

	if currentPage == totalPages {
		// Last page
		singlePageRecordCount = totalRecordCount - (pageSize * (totalPages - 1))
	} else {
		singlePageRecordCount = pageSize
	}

	return &dto.ResponseGetPostsPagination{
		CurrentPage:           currentPage,
		TotalPages:            totalPages,
		SinglePageRecordCount: singlePageRecordCount,
		TotalRecordCount:      totalRecordCount,
	}
}

func (s *service) LikePost(ctx context.Context, postID string, action string, userID string, ChannelID string) *exceptions.Exception {
	log := logger.GetLogInstance(ctx, "Like Post Service")
	post, errPost := s.repo.CheckPostIDValidity(ctx, postID, ChannelID)
	if strings.ToLower(post.Type) == reply && action == bookmark {
		return exceptions.GetExceptionByErrorCode(exceptions.APICallErrorCode)
	}
	if errPost != nil {
		if errors.Is(errPost, gorm.ErrRecordNotFound) {
			return exceptions.GetExceptionByErrorCode(exceptions.PostIdErrorCode)
		}
		log.Errorf("[BookMarkPostService] Not found any posts for postID: %s and ChannelID: %s and error: %v", postID, ChannelID, errPost)
		return exceptions.GetExceptionByErrorCode(exceptions.BadRequestErrorCode)
	}

	if post.Status == deleted {
		log.Errorf("[LikePostService] Cannot like a deleted post for postID: %s and userID: %s and action: %s", postID, userID, action)
		return exceptions.GetExceptionByErrorCode(exceptions.DeletedPostErrorCode)
	}

	_, err := s.repo.ActionSpecificLikePost(ctx, postID, action, userID)
	if err != nil {
		log.Errorf("[LikePostService] Not found any posts for postID: %s and userID: %s and action: %s and error: %v", postID, userID, action, err)
		return err
	}
	return nil
}

func (s *service) DeletePost(ctx context.Context, postID string, userID string) *exceptions.Exception {
	_, err := s.repo.DeleteSpecificPost(ctx, postID, userID)
	log := logger.GetLogInstance(ctx, "Delete Post Service")
	if err != nil {
		log.Errorf("[DeletePostService] Not found any posts to delete for postId: %s with error: %v", postID, err)
		return exceptions.GetExceptionByErrorCode(exceptions.BadRequestErrorCode)
	}
	return nil
}

func (s *service) ReportPost(ctx context.Context, requestBody *model.RequestReportPost, userId string) *exceptions.Exception {

	report := dbModel.Reports{
		PostID:         requestBody.PostID,
		MasterReportID: int64(requestBody.MasterReportID),
		ReportedBy:     userId,
	}

	_, err := s.repo.ReportPostData(ctx, report)
	log := logger.GetLogInstance(ctx, "Report on Post")

	if err != nil {
		log.Error("[ReportPostService] No post found to report for corresponding userId", userId)

		return exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}

	return nil

}

func (s *service) GetRepliesCount(ctx context.Context, postID string) (int64, *exceptions.Exception) {
	count, err := s.repo.GetCommentSpecificReplyCount(ctx, postID)
	log := logger.GetLogInstance(ctx, "GetRepliesCountService")
	if err != nil {
		log.Errorf("[GetRepliesCountService] Couldn't get reply count for corresponding comment id: %s", postID)
		return 0, exceptions.GetExceptionByErrorCode(exceptions.BadRequestErrorCode)
	}
	return count, nil
}

func (s *service) AllRepliesOnPost(ctx context.Context, postId string, userId string, ChannelID string, limit int, currentPage int, sortBy string) (*dto.ResponseAllRepliesOnPost, *exceptions.Exception) {
	log := logger.GetLogInstance(ctx, "AllRepliesOnPost")
	var allReplies []*dto.ResponseAllRepliesOnPostReplies
	comment, err := s.repo.GetPostByPostId(ctx, postId)
	if err != nil {
		log.Error("failed to get post by post id: error: %w", err)
		return nil, exceptions.GetExceptionByErrorCode(exceptions.QueryFailedErrorCode)
	}
	postId = strings.TrimSpace(postId)
	commentPostIdConverted, err := strconv.ParseInt(postId, 10, 32)
	if err != nil {
		log.Printf("Error converting PostId:%s from string to int64: %v", postId, err)
		commentPostIdConverted = 0
	}
	commentLikeStatus, bookmarkStatus, errCommentLikeStatus := s.repo.FetchUserPostSpecificActionValue(ctx, commentPostIdConverted, userId)
	if errCommentLikeStatus != nil {
		commentLikeStatus = false
		log.Errorf("failed to get post by post id: error: %w : %v", err, bookmarkStatus)
	}
	comment.IsLiked = commentLikeStatus
	comment.IsBookmarked = bookmarkStatus
	byteData, err := json.Marshal(comment)
	if err != nil {
		log.Errorf("failed to marshal comment data: %s", err)
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}

	var response dto.ResponseAllRepliesOnPostData

	err = json.Unmarshal(byteData, &response)
	if err != nil {
		log.Errorf("failed to marshal comment data: %s", err)
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}

	// Parse and convert reply CreatedAt and UpdatedAt to Unix timestamps
	parsedCreatedAt, err := time.Parse(time.RFC3339, response.CreatedAt)
	if err == nil {
		response.CreatedAt = fmt.Sprint(parsedCreatedAt.Unix())
	} else {
		log.Errorf("Error parsing CreatedAt time: %v", err)
	}

	parsedUpdatedAt, err2 := time.Parse(time.RFC3339, response.UpdatedAt)
	if err2 == nil {
		response.UpdatedAt = fmt.Sprint(parsedUpdatedAt.Unix())
	} else {
		log.Errorf("Error parsing Updated time: %v", err)
	}

	replies, err := s.repo.GetAllRepliesOnPost(ctx, postId, ChannelID, limit, currentPage, sortBy)
	if err != nil {
		log.Errorf("[AllRepliesOnPostService] failed to fetch replies for postId: %s, got error: %s", postId, err)
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}
	for _, reply := range replies {
		likeStatus, _, errLikeStatus := s.repo.FetchUserPostSpecificActionValue(ctx, reply.ID, userId)
		if errLikeStatus != nil {
			likeStatus = false
			log.Errorf("[AllRepliesOnPostService] failed to fetch like status for replyID: %d, userID: %s, got error: %v", reply.ID, reply.UserID, err)
		}

		parsedCreatedAt, err := time.Parse(time.RFC3339, reply.CreatedAt)
		if err == nil {
			reply.CreatedAt = fmt.Sprint(parsedCreatedAt.Unix())
		} else {
			log.Errorf("[AllRepliesOnPostService]Error parsing reply CreatedAt time: %v", err)
		}

		parsedUpdatedAt, err := time.Parse(time.RFC3339, reply.UpdatedAt)
		if err == nil {
			reply.UpdatedAt = fmt.Sprint(parsedUpdatedAt.Unix())
		} else {
			log.Errorf("[AllRepliesOnPostService]Error parsing reply UpdatedAt time: %v", err)
		}

		userName := reply.FirstName
		if reply.MiddleName != "" {
			userName += " " + reply.MiddleName
		}
		if reply.LastName != "" {
			userName += " " + reply.LastName
		}

		allReplies = append(allReplies, &dto.ResponseAllRepliesOnPostReplies{
			PostID:          fmt.Sprint(reply.ID),
			UserName:        userName,
			UserPhone:       reply.UserPhone,
			ProfileImageURL: reply.ProfileImageUrl,
			Content:         reply.Content,
			Type:            reply.Type,
			LikeCount:       reply.LikeCount,
			Status:          reply.Status,
			IsLiked:         likeStatus,
			UserID:          reply.UserID,
			CreatedAt:       reply.CreatedAt,
			UpdatedAt:       reply.UpdatedAt,
		})
	}
	var dereferencedReplies []dto.ResponseAllRepliesOnPostReplies
	for _, reply := range allReplies {
		dereferencedReplies = append(dereferencedReplies, *reply)
	}

	response.Replies = dereferencedReplies
	repliesCount, errRepliesCount := s.GetRepliesCount(ctx, postId)
	if errRepliesCount != nil {
		repliesCount = int64(len(replies))
		log.Errorf("[AllRepliesOnPostService] failed to fetch replies count for postId: %s, got error: %s", postId, err)
	}

	pagination := NewPagination(int64(currentPage), int64(limit), int64(repliesCount))

	return &dto.ResponseAllRepliesOnPost{
		Code:    APISuccessCode,
		Message: APISuccessMessage,
		Data:    response,
		Pagination: dto.ResponseAllRepliesOnPostPagination{
			CurrentPage:           pagination.CurrentPage,
			TotalPages:            pagination.TotalPages,
			SinglePageRecordCount: pagination.SinglePageRecordCount,
			TotalRecordCount:      pagination.TotalRecordCount,
		},
	}, nil
}

func (s *service) MarkAsRead(ctx context.Context, userID string, ChannelID string) (*dto.ResponseMarkNotificationsAsRead, *exceptions.Exception) {

	isRead, err := s.HasUserReadPost(ctx, userID, ChannelID)
	if err != nil {
		return &dto.ResponseMarkNotificationsAsRead{}, exceptions.GetExceptionByErrorCode(exceptions.QueryFailedErrorCode)
	}

	return &dto.ResponseMarkNotificationsAsRead{
		Code:    APISuccessCode,
		Message: APISuccessMessage,
		Data: dto.ResponseMarkNotificationsAsReadData{
			IsUnread: isRead,
		},
	}, nil
}

func (s *service) HasUserReadPost(ctx context.Context, userID, ChannelID string) (bool, error) {
	log := logger.GetLogInstance(ctx, "HasUserReadChannel-Controller")

	redisKey := fmt.Sprintf("community_comment_unread:%s:%s", userID, ChannelID)

	value, err := s.redisClient.GetKey(ctx, redisKey)
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		log.Errorf("[HasUserReadChannel-Controller] Error retrieving Redis key %s: %v", redisKey, err)
		return false, nil
	}

	if value == "" {
		return false, nil
	}

	return value == "true", nil
}
