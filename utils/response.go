package utils

import (
	"net/http"
	"tangsong-esports/models"

	"github.com/gin-gonic/gin"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, models.Response{
		Code:    models.StatusSuccess,
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应带消息
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, models.Response{
		Code:    models.StatusSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, models.Response{
		Code:    models.StatusError,
		Message: message,
	})
}

// ErrorWithCode 错误响应带状态码
func ErrorWithCode(c *gin.Context, code int, message string) {
	httpStatus := http.StatusBadRequest
	switch code {
	case http.StatusUnauthorized:
		httpStatus = http.StatusUnauthorized
	case models.StatusForbidden:
		httpStatus = http.StatusForbidden
	case models.StatusNotFound:
		httpStatus = http.StatusNotFound
	case models.StatusError:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, models.Response{
		Code:    code,
		Message: message,
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, models.Response{
		Code:    models.StatusForbidden,
		Message: message,
	})
}

// NotFound 未找到
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, models.Response{
		Code:    models.StatusNotFound,
		Message: message,
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, models.Response{
		Code:    models.StatusForbidden,
		Message: message,
	})
}
