package stroke

import (
	"github.com/apacana/apacana-api/biz/config"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/transform"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func CreateStroke(c *gin.Context) {
	var createStrokeForm CreateStrokeForm
	if err := c.ShouldBindJSON(&createStrokeForm); err != nil ||
		len(createStrokeForm.StrokeName) == 0 || len(createStrokeForm.StrokeName) > 24 {
		helper.FormatLogPrint(helper.WARNING, "CreateStroke bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "CreateStroke from: %+v", createStrokeForm)
	userToken := c.GetString(helper.UserToken)
	userInfo, err := mysql.GetUserInfoByToken(c, nil, userToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// insert tourist
			nowTime := time.Now().Format("2006-01-02 15:04:05")
			err := mysql.InsertUserInfo(c, &mysql.UserInfo{
				Token:      userToken,
				UserName:   "",
				PassWord:   "",
				Name:       "",
				Status:     helper.TouristStatus,
				CreateTime: nowTime,
				UpdateTime: nowTime,
			})
			if err != nil {
				helper.FormatLogPrint(helper.ERROR, "CreateStroke InsertUserInfo failed, err: %v, userToken: %v", err, userToken)
				helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
				return
			}
			userInfo, err = mysql.GetUserInfoByToken(c, nil, userToken)
			if err != nil {
				helper.FormatLogPrint(helper.ERROR, "GetUserInfoByToken Tourist GetUserInfoByToken failed, err: %v, userToken: %v", err, userToken)
			}
		} else {
			helper.FormatLogPrint(helper.ERROR, "CreateStroke GetUserInfoByToken failed, err: %v, userToken: %v", err, userToken)
			helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
			return
		}
	}

	if userInfo.Status == helper.TransferredStatus {
		helper.FormatLogPrint(helper.WARNING, "CreateStroke Invalid User, token: %v", userInfo.Token)
		helper.BizResponse(c, http.StatusOK, helper.CodeInvalidUser, nil)
		return
	}

	// quota judge
	strokeList, err := transform.StringToStrokeList(userInfo.Strokes)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateStroke StringToStrokeList failed, strokes: %v", userInfo.Strokes)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	} else if len(strokeList.StrokeList) >= config.StrokeLimit {
		helper.FormatLogPrint(helper.WARNING, "CreateStroke ErrStrokeOutOfLimit, token: %v", userInfo.Token)
		helper.BizResponse(c, http.StatusOK, helper.CodeStrokeOutOfLimit, nil)
		return
	}

	// insert
	strokeToken, err := createUserStroke(c, userInfo, strokeList, createStrokeForm.StrokeName)
	if err != nil {
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
		"stroke_name":  createStrokeForm.StrokeName,
		"stroke_token": strokeToken,
	})
}

func GetStroke(c *gin.Context) {
	strokeToken := c.Param("strokeToken")
	if strokeToken == "" {
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}

	userToken := c.GetString(helper.UserToken)
	userInfo, err := mysql.GetUserInfoByToken(c, nil, userToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetStroke GetUserInfoByToken failed, userInfo: %v", userInfo)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if userInfo.Status == helper.TransferredStatus {
		helper.BizResponse(c, http.StatusOK, helper.CodeInvalidUser, nil)
		return
	}

	strokeInfo, err := mysql.GetStrokeByToken(c, nil, strokeToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetStroke GetStrokeByToken failed, strokeToken: %v", strokeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if strokeInfo.Status == helper.StrokeDeleteStatus {
		helper.FormatLogPrint(helper.LOG, "GetStroke Stroke has deleted, strokeToken: %v", strokeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeMountDeleted, nil)
		return
	}

	if strokeInfo.OwnerID != userInfo.ID {
		helper.FormatLogPrint(helper.LOG, "GetStroke forbidden, strokeToken: %v, userToken: %v", strokeToken, userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}

	routeList, err := transform.CreateFmtRouteList(c, strokeInfo.RoutesList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetStroke CreateFmtRouteList failed, routeList: %v", strokeInfo.RoutesList)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
		"stroke_name":  strokeInfo.StrokeName,
		"stroke_token": strokeToken,
		"route_list":   routeList,
	})
}
