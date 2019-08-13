package middleware

import (
	"fmt"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ApacanaCookieRequire(c *gin.Context) {
	session, err := c.Cookie(helper.ApacanaSession)
	if err != nil {
		helper.BizResponseAndAbort(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}
	elements := strings.Split(session, "-")
	if len(elements) != 4 {
		helper.BizResponseAndAbort(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}

	key := fmt.Sprintf("%s%s%s%s", elements[0], elements[1], elements[2], helper.SessionSalt)
	cipherStr := helper.Md5(key)

	if cipherStr != elements[3] {
		helper.BizResponseAndAbort(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}
	timeStamp, err := strconv.ParseInt(elements[1], 10, 64)
	if err != nil {
		helper.BizResponseAndAbort(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}
	if time.Now().Unix()-timeStamp > helper.YearTime {
		helper.BizResponseAndAbort(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}
	c.Set(helper.UserToken, elements[2])
}
