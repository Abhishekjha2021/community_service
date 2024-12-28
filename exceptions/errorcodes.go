package exceptions

import "net/http"

type ErrorCode string
type ErrorMessage string

const (
	SomethingWentWrongErrorCode   ErrorCode = "MP001"
	QueryFailedErrorCode          ErrorCode = "MPQFE"
	DeletedPostErrorCode          ErrorCode = "MPDPE"
	NoDataFoundErrorCode          ErrorCode = "MP002"
	BadRequestErrorCode           ErrorCode = "MP009"
	APICallErrorCode              ErrorCode = "MPACE"
	PostIdErrorCode               ErrorCode = "MPPIE"
	UnLikeBeforeLikeErrorCode     ErrorCode = "MPUBL"
	AlreadyLikedErrorCode         ErrorCode = "MPALE"
	AlreadyUnLikedErrorCode       ErrorCode = "MPAUE"
	UserIDMissingErrorCode        ErrorCode = "MPUIM"
	QueryParamsIncorrectErrorCode ErrorCode = "MPQPI"
	AccessDeniedErrorCode         ErrorCode = "MPADE"
	WrongChannelIdErrorCode       ErrorCode = "MPWCI"
)

const (
	somethingWentWrong            ErrorMessage = "something went wrong"
	noDataFoundErrorMessage       ErrorMessage = "no data found"
	queryFailedErrorMessage       ErrorMessage = "database query failed"
	badRequestErrorMessage        ErrorMessage = "bad request"
	apiCallErrorMessage           ErrorMessage = "api calling error"
	postIdErrorCode               ErrorMessage = "Post Id not correct in query params"
	unLikeBeforeLikeErrorCode     ErrorMessage = "Cannot unlike before liking the Post"
	alreadyLikedErrorCode         ErrorMessage = "Cannot like already liked post"
	alreadyUnLikedErrorCode       ErrorMessage = "Cannot unLike already unLiked post"
	queryParamsIncorrectErrorCode ErrorMessage = "Query params not correct"
	userIDMissingErrorCode        ErrorMessage = "User id not sent in header"
	accessDeniedErrorMessage      ErrorMessage = "User not authorized to perform action"
	deletedPostErrorMessage       ErrorMessage = "Cannot perform action on a deleted post"
	wrongChannelIdErrorMessage    ErrorMessage = "Given channel-id is wrong"
)

var (
	errorCodeErrorMessageMapping = map[ErrorCode]ErrorMessage{
		SomethingWentWrongErrorCode:   somethingWentWrong,
		NoDataFoundErrorCode:          noDataFoundErrorMessage,
		BadRequestErrorCode:           badRequestErrorMessage,
		QueryFailedErrorCode:          queryFailedErrorMessage,
		APICallErrorCode:              apiCallErrorMessage,
		PostIdErrorCode:               postIdErrorCode,
		UnLikeBeforeLikeErrorCode:     unLikeBeforeLikeErrorCode,
		AlreadyLikedErrorCode:         alreadyLikedErrorCode,
		AlreadyUnLikedErrorCode:       alreadyUnLikedErrorCode,
		UserIDMissingErrorCode:        userIDMissingErrorCode,
		QueryParamsIncorrectErrorCode: queryParamsIncorrectErrorCode,
		AccessDeniedErrorCode:         accessDeniedErrorMessage,
		DeletedPostErrorCode:          deletedPostErrorMessage,
		WrongChannelIdErrorCode:       wrongChannelIdErrorMessage,
	}
)

var (
	errorCodeHttpCodeMapping = map[ErrorCode]int{
		SomethingWentWrongErrorCode:   http.StatusInternalServerError,
		QueryFailedErrorCode:          http.StatusInternalServerError,
		NoDataFoundErrorCode:          http.StatusOK,
		BadRequestErrorCode:           http.StatusBadRequest,
		APICallErrorCode:              http.StatusInternalServerError,
		PostIdErrorCode:               http.StatusBadRequest,
		UnLikeBeforeLikeErrorCode:     http.StatusBadRequest,
		AlreadyLikedErrorCode:         http.StatusBadRequest,
		AlreadyUnLikedErrorCode:       http.StatusBadRequest,
		UserIDMissingErrorCode:        http.StatusBadRequest,
		QueryParamsIncorrectErrorCode: http.StatusBadRequest,
		AccessDeniedErrorCode:         http.StatusUnauthorized,
		DeletedPostErrorCode:          http.StatusBadRequest,
		WrongChannelIdErrorCode:       http.StatusBadRequest,
	}
)

func (c ErrorCode) String() string {
	return string(c)
}

func (c ErrorCode) GetErrorMessage() ErrorMessage {
	if errMsg, ok := errorCodeErrorMessageMapping[c]; ok {
		return errMsg
	}

	return somethingWentWrong
}

func (m ErrorMessage) String() string {
	return string(m)
}
