package stroke

import (
	"github.com/apacana/apacana-api/biz/config"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func CreateStroke(c *gin.Context) {
	var createStrokeForm CreateStrokeForm
	if err := c.ShouldBindJSON(&createStrokeForm); err != nil || len(createStrokeForm.StrokeName) > 40 {
		helper.FormatLogPrint(helper.WARNING, "CreateStroke bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
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
	// quota judge
	strokeList := &helper.StrokeList{
		StrokeList: make([]int64, 0),
	}
	if len(userInfo.Strokes) != 0 {
		strokeList = helper.StringToStrokeList(userInfo.Strokes)
		if strokeList == nil {
			helper.FormatLogPrint(helper.ERROR, "CreateStroke StringToStrokeList failed, strokes: %v", userInfo.Strokes)
			helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
			return
		} else if len(strokeList.StrokeList) >= config.StrokeLimit {
			helper.BizResponse(c, http.StatusOK, helper.CodeStrokeOutOfLimit, nil)
			return
		}
	}

	// insert
	strokeName, err := createUserStroke(c, userInfo, strokeList, createStrokeForm.StrokeName)
	if err != nil {
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{"stroke_name": strokeName})
}
