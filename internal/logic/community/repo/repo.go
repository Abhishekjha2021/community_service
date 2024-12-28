package repo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
	"github.com/Abhishekjha321/community_service/log"
	"github.com/Abhishekjha321/community_service/internal/common"
	"github.com/Abhishekjha321/community_service/internal/logic/community/model"
	dbModel "github.com/Abhishekjha321/community_service/pkg/store/db/model"

	"github.com/Abhishekjha321/community_service/exceptions"
	"github.com/Abhishekjha321/community_service/pkg/store/db"
	"gorm.io/gorm"
)

const (
	postsTable       = "posts"
	userActionsTable = "user_actions"
	like             = common.Like_Action
	unlike           = "unlike"
	bookmark         = common.Bookmark_Action
	reportTable      = "reports"
	userDetailsTable = "user_details"
	allReplyTable    = "replies"
	comment          = "COMMENT"
)

type repo struct {
	db *db.Store
}

func NewRepo(db *db.Store) model.Repo {

	return &repo{
		db: db,
	}
}

// func (r *repo) GetReportCategoriesData(ctx context.Context) ([]model.ReportResponse, error) {
// 	log := logger.GetLogInstance(ctx, "GetReportCategoriesData")
// 	var result []model.ReportResponse
// 	err := r.db.MasterDB.WithContext(ctx).Raw(
// 		`Select parent_reports.id as parent_id,
// 		parent_reports.title as parent_title,
// 		parent_reports.topic as parent_topic,
// 		parent_reports.subtitle as parent_subtitle,
// 		child_reports.id as child_id,
// 		child_reports.title as child_title,
// 		child_reports.topic as child_topic,
// 		child_reports.subtitle as child_subtitle
// 	from
// 		master_reports as parent_reports
// 	left join master_reports as child_reports on
// 		child_reports.parent_id = parent_reports.id
// 	where
// 		parent_reports.id in (select id from master_reports where parent_id = 0)`).Find(&result).Error
// 	if err != nil {
// 		log.Errorf("[GetReportCategoriesData] Error while retreiving data from db with error: %v", err)
// 		return nil, err
// 	}

// 	return result, nil
// }

// func (r *repo) GetEventPostEngagementData(ctx context.Context, channelID string) (*model.LikeCommentCountOnChannel, error) {
// 	log := logger.GetLogInstance(ctx, "GetEventPostEngagementData")
// 	var result model.LikeCommentCountOnChannel
// 	err := r.db.MasterDB.WithContext(ctx).Table(postsTable).Select("COUNT(*) as comment_count, SUM(like_count) as like_count").Where("channel_id = ?", channelID).Find(&result).Error
// 	if err != nil {
// 		log.Errorf("[GetEventPostEngagementDataRepo] Error while retreiving data from db with error: %v for channelID: %s", err, channelID)
// 		return nil, err
// 	}

// 	return &result, nil
// }

func (r *repo) GetUserIDByPostID(ctx context.Context, postID int64) (string, error) {
	log := logger.GetLogInstance(ctx, "GetUserIDByPostID-repo")

	var userID string
	db := r.db.MasterDB.WithContext(ctx).Table(postsTable).Select("user_id").Where("id = ?", postID).Scan(&userID)
	if db.Error != nil {
		log.Errorf("[GetUserIDByPostIDRepo] error while fetching user_id from db: %+v", db.Error.Error())
		return userID, nil
	}

	return userID, nil
}

func (r *repo) GetUserDetailsByUserId(ctx context.Context, userId string) (*dbModel.UserDetails, error) {
	log := logger.GetLogInstance(ctx, "GetUserDetails-repo")

	var userDetails dbModel.UserDetails
	db := r.db.MasterDB.WithContext(ctx).Table(userDetailsTable).Where("user_id = ?", userId).First(&userDetails)
	if db.Error != nil {
		if errors.Is(db.Error, gorm.ErrRecordNotFound) {
			log.Infof("[GetUserDetailsByUserIdRepo] no user details found for user_id: %s", userId)
			return &userDetails, nil
		}
		log.Errorf("[GetUserDetailsByUserIdRepo] error while fetching user details from db: %+v", db.Error.Error())
		return nil, fmt.Errorf("getUserDetails query failed: %w", db.Error)
	}

	return &userDetails, nil
}

func (r *repo) InsertPostData(ctx context.Context, postData dbModel.Post) (*dbModel.Post, error) {

	log := logger.GetLogInstance(ctx, "InsertPostData-repo")

	db := r.db.MasterDB.WithContext(ctx).Table(postsTable).Create(&postData)
	if db.Error != nil {
		log.Errorf("[InsertPostDataRepo] error while creating data in db: %+v", db.Error.Error())
		return nil, fmt.Errorf("createPost query failed: %w", db.Error)
	}
	if db.RowsAffected < 1 {
		log.Errorf("[InsertPostDataRepo] failed to insert post data in db")
		return nil, errors.New("[InsertPostDataRepo] failed to insert post data in db")
	}
	return &postData, nil

}

func (r *repo) GetEventPostsCount(ctx context.Context, channelID string) (int64, error) {
	log := logger.GetLogInstance(ctx, "GetEventPostsCount")
	var count int64
	db := r.db.MasterDB.WithContext(ctx).Table(postsTable).Where("channel_id = ? AND type = ?", channelID, comment).Count(&count)
	if db.Error != nil {
		log.Errorf("[GetEventPostsCount] Error while fetching count of posts for channel id: %s from db with error: %+v", channelID, db.Error)
		return 0, db.Error
	}
	return count, nil
}

func (r *repo) GetUserSpecificPostsCount(ctx context.Context, channelId string, userId string) (int, error) {
	var count int64
	log := logger.GetLogInstance(ctx, "GetUserSpecificPostsCount")
	db := r.db.MasterDB.WithContext(ctx).Table(postsTable).Where("channel_id = ? AND type = ? AND user_id = ?", channelId, comment, userId).Count(&count)
	if db.Error != nil {
		log.Errorf("[GetUserSpecificPostsCount] Error while fetching count of posts for channel id: %s  and userId : %s from db with error: %+v", channelId, userId, db.Error)
		return 0, db.Error
	}
	return int(count), nil
}

func (r *repo) GetUserSpecificEventPosts(ctx context.Context, channelID string, limit int, currentPage int, userId string, offset int) ([]common.Post, error) {
	var posts []common.Post
	db := r.db.MasterDB.WithContext(ctx).Raw(`
	SELECT 
		id,
		channel_id,
		user_id,
		content,
		type,
		parent_id,
		like_count,
		status,
		created_at,
		updated_at,
		deleted_at
		FROM posts 
		WHERE user_id = ? AND channel_id = ? AND type = 'COMMENT'
		ORDER BY
			CASE 
				WHEN status = 'DELETED' THEN 1
				ELSE 0
			END,
			updated_at desc
			LIMIT ? OFFSET ?
	`, userId, channelID, limit, offset)
	result := db.Find(&posts)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}
	return posts, nil
}

func (r *repo) GetRelevantReplies(ctx context.Context, commentIds []int64, sortBy string, bookMarksOnly bool) ([]common.Post, error) {
	log := logger.GetLogInstance(ctx, "Get Relevant Replies")
	var replies []common.Post
	if bookMarksOnly {
		return replies, nil
	}
	baseQuery := `
		WITH RankedPosts AS (
			SELECT 
				p.*, 
				ROW_NUMBER() OVER (
					PARTITION BY parent_id 
					ORDER BY`
	if sortBy == common.IDEAS_BASED_FLOW {
		baseQuery += `
		CASE 
			WHEN is_pinned = true THEN 0 
			ELSE 1 
		END,`
	}
	db := r.db.MasterDB.WithContext(ctx).Raw(`
					`+baseQuery+`  
					CASE 
						WHEN status = 'DELETED' THEN 1 
						ELSE 0 
					END, 
					updated_at DESC
			) AS row_num
		FROM posts AS p
		WHERE parent_id IN (?)
	)
	SELECT *
	FROM RankedPosts
	WHERE row_num <= 3`, commentIds).Scan(&replies)
	if err := db.Error; err != nil {
		log.Errorf("[GetRelevantReplies] Unable to fetch relevant replies from db")
		return replies, err
	}
	return replies, nil
}

func (r *repo) GetBookMarkedPosts(ctx context.Context, channelID string, limit int, currentPage int, userId string, offset int, sortBy string) ([]common.Post, error) {
	var posts []common.Post
	if sortBy == common.USER_BASED_FLOW {
		return posts, nil
	}
	db := r.db.MasterDB.WithContext(ctx).Raw(`
			select
			ua.post_id as id,
			p.channel_id,
			ua.user_id,
			p.content,
			p.type,
			p.parent_id,
			p.like_count,
			p.bookmark_count,
			p.status,
			ua.created_at,
			ua.updated_at,
			p.deleted_at
			from
				posts p
			left join user_actions ua 
			on
				p.id = ua.post_id
			where
				ua.action = 'bookmark'
				and ua.user_id = ?
				and p.channel_id = ?
				and ua.value = true
			order by
				ua.updated_at desc
			limit ? offset ?	
	`, userId, channelID, limit, offset)
	result := db.Find(&posts)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}
	return posts, nil
}

func (r *repo) GetEventPosts(ctx context.Context, channelID string, limit int, currentPage int, userId string, offset int, sortBy string) ([]common.Post, error) {
	var posts []common.Post
	whereClause := `WHERE user_id != ? AND channel_id = ? AND type = 'COMMENT'`
	orderByClause := `ORDER BY
			CASE 
				WHEN status = 'DELETED' THEN 1
				ELSE 0
			END,
			like_count DESC,
			updated_at desc`
	querParams := []interface{}{userId, channelID, limit, offset}

	if sortBy == common.IDEAS_BASED_FLOW {
		whereClause = `WHERE channel_id = ? AND type = 'COMMENT'`
		orderByClause = `ORDER BY
			CASE 
				WHEN is_pinned = true THEN 0
				ELSE 1
			END,
			CASE 
				WHEN status = 'DELETED' THEN 1
				ELSE 0
			END,
			updated_at desc`
		querParams = []interface{}{channelID, limit, offset}
	}

	db := r.db.MasterDB.WithContext(ctx).Raw(`
	SELECT 
		id AS id,
		channel_id AS channel_id,
		user_id AS user_id,
		content AS content,
		type AS type,
		is_pinned,
		parent_id AS parent_id,
		like_count AS like_count,
		bookmark_count AS bookmark_count,
		status AS status,
		created_at AS created_at,
		updated_at AS updated_at,
		deleted_at AS deleted_at
		FROM posts
		`+whereClause+`
		`+orderByClause+`
		LIMIT ? OFFSET ?
	`, querParams...)
	result := db.Find(&posts)
	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}
	return posts, nil
}

func (r *repo) FetchUserPostSpecificActionValue(ctx context.Context, postID int64, userID string) (bool, bool, error) {
	var userActions []model.UserAction
	result := r.db.MasterDB.WithContext(ctx).
		Table(userActionsTable).
		Select("action, value").
		Where("post_id = ? AND user_id = ? AND action IN (?, ?)", postID, userID, "like", "bookmark").
		Find(&userActions)

	if result.Error != nil {
		return false, false, fmt.Errorf("failed to fetch action values: %w", result.Error)
	}

	likeValue, bookmarkValue := false, false

	for _, action := range userActions {
		switch action.Action {
		case like:
			likeValue = action.Value
		case bookmark:
			bookmarkValue = action.Value
		default:
			likeValue = action.Value
		}
	}

	return likeValue, bookmarkValue, nil
}

func (r *repo) GetUserDetailsForPostID(ctx context.Context, postIDs []int) ([]model.PostUserDetails, error) {
	var data []model.PostUserDetails
	result := r.db.MasterDB.WithContext(ctx).Raw(`
	select p.id,p.user_id, ud.first_name as first_name , ud.user_phone, ud.profile_image_url , ua."action" , ua.value 
	from posts p 
	left join user_details ud 
	on p.user_id  = ud.user_id
	left join user_actions ua 
	on p.id = ua.post_id  and p.user_id  = ua.user_id 
	where p.id in(?)`, postIDs).Find(&data)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch user details: %w", result.Error)
	}

	return data, nil
}

func (r *repo) CheckPostIDValidity(ctx context.Context, postID string, channelId string) (common.Post, error) {
	var post common.Post
	log := logger.GetLogInstance(ctx, "Check Post ID Validity")
	convertedPostID, err := strconv.Atoi(postID)
	if err != nil {
		log.Errorf("[CheckPostIDValidity] couldn't convert postId: %s to integer", postID)
		return post, exceptions.GetExceptionByErrorCode(exceptions.PostIdErrorCode)
	}
	if len(channelId) != 0 {
		errConfirmPost := r.db.MasterDB.WithContext(ctx).Table(postsTable).
			Where("id = ? AND channel_id = ?", convertedPostID, channelId).
			Scan(&post).Error
		if errConfirmPost != nil {
			return post, errConfirmPost
		}
		return post, nil
	}
	errFetchPost := r.db.MasterDB.WithContext(ctx).Table(postsTable).
		Where("id = ?", convertedPostID).
		Scan(&post).Error

	if errFetchPost != nil || post.ID == 0 {
		if post.ID == 0 {
			return post, gorm.ErrRecordNotFound
		}
		return post, errFetchPost
	}
	return post, nil
}

func (r *repo) UpdateRequiredActionInPostsTable(ctx context.Context, postID string, updateExpr string, userID string, actionName string) *exceptions.Exception {
	log := logger.GetLogInstance(ctx, "UpdateRequiredActionInPostsTable")
	resPosts := r.db.MasterDB.WithContext(ctx).Table(postsTable).Where("id = ?", postID).Update(GetFieldNameBasedOnAction(actionName), gorm.Expr(updateExpr, 1))
	if resPosts.Error != nil {
		log.Errorf("[UpdateRequiredActionInPostsTable] Unable to update like count for user with userID: %s , postID: %s and action: %s with error: %v", userID, postID, actionName, resPosts.Error)
		return exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}
	return nil
}

func GetFieldNameBasedOnAction(actionName string) string {
	switch actionName {
	case like:
		return "like_count"
	case bookmark:
		return "bookmark_count"
	default:
		return ""
	}
}

func (r *repo) ActionSpecificLikePost(ctx context.Context, postID string, action string, userID string) (string, *exceptions.Exception) {
	log := logger.GetLogInstance(ctx, "Action Specific Like Post")
	convertedPostID, err := strconv.Atoi(postID)
	if err != nil {
		log.Errorf("[ActionSpecificLikePost] couldn't convert postId: %s to integer", postID)
		return "", exceptions.GetExceptionByErrorCode(exceptions.PostIdErrorCode)
	}
	actionName := action
	if actionName == like || actionName == unlike {
		actionName = like
	}
	var userActions []*dbModel.UserActions
	checkAction := r.db.MasterDB.WithContext(ctx).Table(userActionsTable).Where("post_id = ? AND user_id = ? AND action = ?", postID, userID, actionName).Find(&userActions)
	if checkAction.Error != nil {
		log.Errorf("[ActionSpecificLikePost] unable to fetch like value for postID: %s with error: %v", postID, checkAction.Error)
		return "", exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
	}
	var updateExpr string
	var newValue bool
	postColumnToUpdate := GetFieldNameBasedOnAction(actionName)
	if len(userActions) == 0 || !*userActions[0].Value {
		updateExpr = postColumnToUpdate + " + ?"
		newValue = true
	} else {
		updateExpr = postColumnToUpdate + " - ?"
		newValue = false
	}
	errUpdatePostTable := r.UpdateRequiredActionInPostsTable(ctx, postID, updateExpr, userID, actionName)
	if errUpdatePostTable != nil {
		return "", errUpdatePostTable
	}
	if len(userActions) > 0 {
		updateLike := r.db.MasterDB.WithContext(ctx).Table(userActionsTable).Where("post_id = ? AND user_id = ? AND action = ?", postID, userID, actionName).Updates(&dbModel.UserActions{
			Value: &newValue,
		})
		if updateLike.Error != nil {
			log.Errorf("[ActionSpecificLikePost] Unable to update user actions Action field bool value for userID: %s with postID: %s and action: %s with error: %v", userID, postID, actionName, updateLike.Error)
			return "", exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
		}
	} else {
		userActionCreated := r.db.MasterDB.WithContext(ctx).Table(userActionsTable).Create(&dbModel.UserActions{
			UserID: userID,
			PostID: int64(convertedPostID),
			Action: actionName,
			Value:  &newValue,
		})
		if userActionCreated.Error != nil {
			log.Errorf("[ActionSpecificLikePost] Unable to create user action for userID: %s with postID: %s and action: %s with error: %v", userID, postID, actionName, userActionCreated.Error)
			return "", exceptions.GetExceptionByErrorCode(exceptions.SomethingWentWrongErrorCode)
		}
	}
	return "", nil
}

func (r *repo) DeleteSpecificPost(ctx context.Context, postID string, userID string) (string, error) {
	log := logger.GetLogInstance(ctx, "Delete Specific Post")
	var postSpecificUserID string
	resUserId := r.db.MasterDB.WithContext(ctx).Table(postsTable).Select("user_id").Where("id = ?", postID).Scan(&postSpecificUserID)
	if resUserId.Error != nil {
		log.Errorf("[DeleteSpecificPost] Error while getting userID for postID: %s with error: %v", postID, resUserId.Error)
		return "", exceptions.GetExceptionByErrorCode(exceptions.NoDataFoundErrorCode)
	}
	if postSpecificUserID != userID {
		return "", exceptions.GetExceptionByErrorCode(exceptions.AccessDeniedErrorCode)
	}
	res := r.db.MasterDB.WithContext(ctx).Table(postsTable).Where("id = ?", postID).Updates(&dbModel.Post{
		Content:   "This comment was deleted by the post author.",
		Status:    "DELETED",
		DeletedAt: time.Now(),
	})
	if res.Error != nil {
		log.Errorf("[DeleteSpecificPost] Error while deleting post with postID: %s with error: %v", postID, res.Error)
		return "", exceptions.GetExceptionByErrorCode(exceptions.NoDataFoundErrorCode)
	}
	return "", nil
}

func (r *repo) ReportPostData(ctx context.Context, reportData dbModel.Reports) (*dbModel.Reports, error) {
	log := logger.GetLogInstance(ctx, "ReportPostData-repo")

	db := r.db.MasterDB.WithContext(ctx).Table(reportTable).Create(&reportData)
	if db.Error != nil {
		log.Errorf("[ReportPostDataRepo] error while creating report data in db: %+v", db.Error.Error())
		return nil, fmt.Errorf("reportPost query failed: %w", db.Error)
	}
	if db.RowsAffected < 1 {
		log.Errorf("[ReportPostDataRepo] failed to insert report data in db")
		return nil, errors.New("[ReportPostDataRepo] failed to insert report data in db")
	}
	return &reportData, nil

}

func (r *repo) PopulateUserInfoTable(ctx context.Context, userInfo common.UserInfo) error {
	var db *gorm.DB
	if userInfo.IsNewUser {
		db = r.db.MasterDB.WithContext(ctx).Table(userDetailsTable).Create(&dbModel.UserDetails{
			UserID:          userInfo.UserID,
			ProfileImageUrl: userInfo.ProfileImageUrl,
			FirstName:       userInfo.FirstName,
			MiddleName:      userInfo.MiddleName,
			LastName:        userInfo.LastName,
			Email:           userInfo.Email,
			UserPhone:       userInfo.UserPhone,
			UserName:        userInfo.UserName,
		})
	} else {
		db = r.db.MasterDB.WithContext(ctx).Table(userDetailsTable).Where("user_id = ?", userInfo.UserID).Updates(&dbModel.UserDetails{
			ProfileImageUrl: userInfo.ProfileImageUrl,
			FirstName:       userInfo.FirstName,
			MiddleName:      userInfo.MiddleName,
			LastName:        userInfo.LastName,
			Email:           userInfo.Email,
			UserPhone:       userInfo.UserPhone,
			UserName:        userInfo.UserName,
		})
	}
	if db.Error != nil {
		return fmt.Errorf("[PopulateUserInfoTableRepo] user info query failed: %w", db.Error)
	}

	return nil
}

func (r *repo) GetCommentSpecificReplyCount(ctx context.Context, postID string) (int64, error) {
	log := logger.GetLogInstance(ctx, "GetCommentSpecificReplyCount")
	var count int64
	db := r.db.MasterDB.WithContext(ctx).Table(postsTable).Where("parent_id = ?", postID).Count(&count)
	if db.Error != nil {
		log.Errorf("[GetCommentSpecificReplyCount] Error while fetching count of replies for comment id: %s from db with error: %+v", postID, db.Error)
		return 0, db.Error
	}
	return count, nil
}

func (r *repo) GetPostByPostId(ctx context.Context, postId string) (*common.AllRepliesPost, error) {
	log := logger.GetLogInstance(ctx, "GetPostByPostId-repo")

	var post *common.AllRepliesPost
	db := r.db.MasterDB.WithContext(ctx).
		Table("posts as p").
		Select("p.id as post_id, p.content as content, p.type as type, p.like_count as like_count, p.status as status, p.created_at as created_at,p.updated_at as updated_at, p.user_id as user_id, u.first_name as user_name, u.profile_image_url as profile_image_url, u.user_phone, p.bookmark_count as bookmark_count").
		Joins("left join user_details u on p.user_id = u.user_id").
		Where("p.id = ?", postId).
		Scan(&post)
	if db.Error != nil {
		log.Errorf("[GetPostByPostId] error while executing GetPostByPostId query: %+v", db.Error)
		return nil, fmt.Errorf("GetPostByPostId query failed: %w", db.Error)
	}
	return post, nil
}

func (r *repo) GetAllRepliesOnPost(ctx context.Context, postId string, channelID string, limit int, currentPage int, sortBy string) ([]model.ReplyPost, error) {
	log := logger.GetLogInstance(ctx, "GetAllRepliesOnPost-repo")
	orderClause := `
    CASE 
        WHEN p.status = 'DELETED' THEN 1
        ELSE 0
    END, 
    p.updated_at DESC`
	if sortBy == common.IDEAS_BASED_FLOW {
		orderClause = `
		CASE 
		WHEN p.is_pinned = true THEN 0 
		ELSE 1 
		END, ` + orderClause
	}
	var results []model.ReplyPost
	offset := (currentPage - 1) * limit
	db := r.db.MasterDB.WithContext(ctx).
		Table("posts as p").
		Select("p.id as id, p.content as content, p.type as type, p.like_count as like_count, p.status as status, p.created_at as created_at, p.updated_at as updated_at, u.first_name as first_name, u.middle_name as middle_name, u.last_name as last_name, u.profile_image_url as profile_image_url, p.user_id as user_id, u.user_phone").
		Joins("left join user_details u on p.user_id = u.user_id").
		Where("p.parent_id = ?", postId).
		Order(orderClause).
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&results)

	if db.Error != nil {
		log.Errorf("error while executing GetAllRepliesOnPost query: %+v", db.Error.Error())
		return nil, fmt.Errorf("GetAllRepliesOnPost query failed: %w", db.Error)
	}

	return results, nil
}
