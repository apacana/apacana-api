package middleware

import (
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApacanaCookieRequire(c *gin.Context) {
	session, err := c.Cookie(helper.ApacanaSession)
	if err != nil {
		helper.BizResponseAndAbort(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}
	userToken, err := helper.IsValidCookie(session)
	if err != nil {
		helper.FormatLogPrint(helper.WARNING, "ApacanaCookieRequire IsValidCookie failed, err: %v", err)
		helper.BizResponseAndAbort(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}
	c.Set(helper.UserToken, userToken)
}
