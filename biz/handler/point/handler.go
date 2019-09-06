package point

import (
	"github.com/apacana/apacana-api/biz/config"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/handler/stroke"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/transform"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func AddPoint(c *gin.Context) {
	var addPointForm AddPointForm
	if err := c.ShouldBindJSON(&addPointForm); err != nil ||
		len(addPointForm.Text) > 100 || len(addPointForm.Center) > 50 ||
		len(addPointForm.PlaceName) > 300 || len(addPointForm.PointID) > 50 {
		helper.FormatLogPrint(helper.WARNING, "AddPoint bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	pointType, err := helper.GetPointTypeByName(addPointForm.PointType)
	if err != nil {
		helper.FormatLogPrint(helper.WARNING, "AddPoint GetPointTypeByName failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "AddPoint from: %+v", addPointForm)

	// judge user
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
				helper.FormatLogPrint(helper.ERROR, "AddPoint InsertUserInfo failed, err: %v, userToken: %v", err, userToken)
				helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
				return
			}
			userInfo, err = mysql.GetUserInfoByToken(c, nil, userToken)
			if err != nil {
				helper.FormatLogPrint(helper.ERROR, "AddPoint GetUserInfoByToken failed, err: %v, userToken: %v", err, userToken)
			}
		} else {
			helper.FormatLogPrint(helper.ERROR, "AddPoint GetUserInfoByToken failed, err: %v, userToken: %v", err, userToken)
			helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
			return
		}
	}

	if userInfo.Status == helper.TransferredStatus {
		helper.FormatLogPrint(helper.WARNING, "AddPoint Invalid User, token: %v", userInfo.Token)
		helper.BizResponse(c, http.StatusOK, helper.CodeInvalidUser, nil)
		return
	}

	userStrokeList, err := transform.StringToStrokeList(userInfo.Strokes)
	defaultStrokeID := userStrokeList.DefaultStroke
	if defaultStrokeID == 0 {
		strokeID, _, _, err := stroke.CreateUserStroke(c, userInfo, userStrokeList, config.DefaultStrokeName)
		if err != nil {
			helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
			return
		}
		defaultStrokeID = strokeID
	} else {
		// 唯一性校验
		allow, err := mysql.AllowAddPoint(c, nil, addPointForm.PointID, pointType, defaultStrokeID)
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "AddPoint AllowAddPoint failed, err: %v", err)
			helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
			return
		}
		if !allow {
			helper.BizResponse(c, http.StatusOK, helper.CodePointExist, nil)
			return
		}
	}

	// quota judge
	strokeInfos, err := mysql.MGetStrokeByID(c, nil, []int64{defaultStrokeID})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "AddPoint MGetStrokeByID failed, err: %v, id: %d", err, defaultStrokeID)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	strokeInfo := strokeInfos[0]
	pointList, err := transform.StringToPointList(strokeInfo.PointsList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "AddPoint StringToPointList failed, pointList: %v", strokeInfo.PointsList)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	} else if len(pointList.PointList) >= config.PointLimit {
		helper.FormatLogPrint(helper.WARNING, "AddPoint OutOfLimit, strokeToken: %v", strokeInfo.StrokeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodePointOutOfLimit, nil)
		return
	}

	// insert point
	outPut, err := addStrokePointList(c, strokeInfo, pointList, addPointForm, pointType)
	if err != nil {
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, outPut)
}
