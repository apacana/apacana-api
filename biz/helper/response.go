package helper

import (
	"github.com/gin-gonic/gin"
)

type GinResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func UpdateReportStatus(c *gin.Context, httpCode int, bizCode int32) {
	c.Set(ReportHttpCode, httpCode)
	c.Set(ReportBizCode, bizCode)
}

func NewGinResponse(code int32, result interface{}) *GinResponse {
	msg := "Success"
	if code == 1 {
		msg = "Failed"
	} else if code == 4 {
		msg = "Forbidden"
	}
	return &GinResponse{
		Code:    code,
		Message: msg,
		Data:    result,
	}
}

func BizResponse(c *gin.Context, httpCode int, bizCode int32, data interface{}) {
	UpdateReportStatus(c, httpCode, bizCode)

	resp := NewGinResponse(bizCode, data)
	c.JSON(httpCode, resp)
}

func BizResponseAndAbort(c *gin.Context, httpCode int, bizCode int32, data interface{}) {
	UpdateReportStatus(c, httpCode, bizCode)

	resp := NewGinResponse(bizCode, data)
	c.JSON(httpCode, resp)
	c.Abort()
}

func AbortWithBizResponse(c *gin.Context, httpCode int, bizCode int32, data interface{}) {
	c.Abort()
	BizResponse(c, httpCode, bizCode, data)
}

func BizStatus(c *gin.Context, httpCode int) {
	bizCode := CodeSuccess
	if httpCode >= 500 {
		bizCode = CodeFailed
	}
	UpdateReportStatus(c, httpCode, bizCode)

	c.Status(httpCode)
}
