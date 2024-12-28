package api

import (
	"net/http"

	"github.com/Abhishekjha321/community_service/exceptions"
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	// code
	Code string `json:"code,omitempty"`

	// message
	Message string `json:"message,omitempty"`
}

type SuccessResp struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func SendApiResponseV1(ctx *gin.Context, apiResp interface{}, apiErr *exceptions.Exception) {
	if apiErr != nil {
		ctx.JSON(apiErr.HttpCode, &ErrorResponse{
			Code:    apiErr.ErrorCode.String(),
			Message: apiErr.ErrorMessage.String(),
		})
		return
	}

	if apiResp != nil {
		ctx.JSON(http.StatusOK, apiResp)
		return
	}

	ctx.JSON(http.StatusOK, SuccessResp{
		Code:    "00000",
		Message: "SUCCESS",
	})
}
