package route

import (
	"github.com/apacana/apacana-api/biz/config"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/out"
	"github.com/apacana/apacana-api/biz/transform"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func CreateRoute(c *gin.Context) {
	var createRouteForm CreateRouteForm
	if err := c.ShouldBindJSON(&createRouteForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "CreateRoute bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "CreateRoute from: %+v", createRouteForm)
	routeName := config.DefaultRouteName
	if createRouteForm.RouteName != nil && len(*createRouteForm.RouteName) < 24 && len(*createRouteForm.RouteName) > 0 {
		routeName = *createRouteForm.RouteName
	}

	// judge user
	userToken := c.GetString(helper.UserToken)
	userInfo, err := mysql.GetUserInfoByToken(c, nil, userToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateRoute GetUserInfoByToken failed, err: %v, userToken: %v", err, userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	if userInfo.Status == helper.TransferredStatus {
		helper.FormatLogPrint(helper.WARNING, "CreateRoute Invalid User, token: %v", userInfo.Token)
		helper.BizResponse(c, http.StatusOK, helper.CodeInvalidUser, nil)
		return
	}

	// judge stroke
	strokeInfo, err := mysql.GetStrokeByToken(c, nil, createRouteForm.StrokeToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateRoute GetStrokeByToken failed, err: %v, userToken: %v", err, userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	if strokeInfo.Status == helper.StrokeDeleteStatus {
		helper.FormatLogPrint(helper.LOG, "CreateRoute Stroke has deleted, strokeToken: %v", createRouteForm.StrokeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeMountDeleted, nil)
		return
	}

	if strokeInfo.OwnerID != userInfo.ID {
		helper.FormatLogPrint(helper.LOG, "CreateRoute forbidden, strokeToken: %v, userToken: %v", createRouteForm.StrokeToken, userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}

	// route quota judge
	routeList, err := transform.StringToRouteList(strokeInfo.RoutesList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateRoute StringToRouteList failed, routeList: %v", strokeInfo.RoutesList)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	} else if len(routeList.RouteList) >= config.RouteLimit {
		helper.FormatLogPrint(helper.WARNING, "CreateRoute CodeRouteOutOfLimit, strokeToken: %v", strokeInfo.StrokeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeRouteOutOfLimit, nil)
		return
	}

	// insert route
	routeToken, updateTime, err := createStrokeRoute(c, strokeInfo, routeList, routeName)
	if err != nil {
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, &out.RouteInfoOut{
		RouteToken: routeToken,
		RouteName:  routeName,
		UpdateTime: updateTime,
	})
}

func AddRoutePoint(c *gin.Context) {
	var addRoutePointForm AddRoutePointForm
	if err := c.ShouldBindJSON(&addRoutePointForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "AddRoutePoint bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "AddRoutePoint from: %+v", addRoutePointForm)

	// judge user
	userToken := c.GetString(helper.UserToken)
	userInfo, err := mysql.GetUserInfoByToken(c, nil, userToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "AddRoutePoint GetUserInfoByToken failed, err: %v, userToken: %v", err, userToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	if userInfo.Status == helper.TransferredStatus {
		helper.FormatLogPrint(helper.WARNING, "AddRoutePoint Invalid User, token: %v", userInfo.Token)
		helper.BizResponse(c, http.StatusOK, helper.CodeInvalidUser, nil)
		return
	}

	// judge route and point
	routeInfo, err := mysql.GetRouteByToken(c, nil, addRoutePointForm.RouteToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "AddRoutePoint GetRouteByToken failed, err: %v, routeToken: %v", err, addRoutePointForm.RouteToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if routeInfo.Status == helper.RouteDeleteStatus {
		helper.FormatLogPrint(helper.LOG, "AddRoutePoint Route has deleted, routeToken: %v", addRoutePointForm.RouteToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeMountDeleted, nil)
		return
	}

	// route鉴权
	if routeInfo.OwnerId != userInfo.ID {
		helper.BizResponse(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}

	pointInfo, err := mysql.GetPointByToken(c, nil, addRoutePointForm.PointToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "AddRoutePoint GetPointByToken failed, err: %v, routeToken: %v", err, addRoutePointForm.PointToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	if routeInfo.StrokeID != pointInfo.StrokeID {
		helper.FormatLogPrint(helper.WARNING, "AddRoutePoint route and point don't have same stroke.")
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	strokeList, err := transform.StringToStrokeList(userInfo.Strokes)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "AddRoutePoint StringToStrokeList failed, strokes: %v", userInfo.Strokes)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if strokeList.DefaultStroke != routeInfo.StrokeID {
		helper.FormatLogPrint(helper.WARNING, "AddRoutePoint want to change with not default stroke")
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	// update route
	routePointList, err := transform.StringToRoutePointList(routeInfo.PointsList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "AddRoutePoint StringToRoutePointList failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	routePointList.PointList = append(routePointList.PointList, pointInfo.ID)
	routePointList.DirectionList = append(routePointList.DirectionList, 0)

	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = mysql.UpdateRouteByToken(c, nil, routeInfo.RouteToken, map[string]interface{}{
		"update_time": nowTime,
		"points_list": *transform.PackRoutePointList(routePointList),
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "AddRoutePoint UpdateRouteByToken failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}

	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, nil)
}

func GetRoute(c *gin.Context) {
	routeToken := c.Param("routeToken")
	if routeToken == "" {
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}

	userToken := c.GetString(helper.UserToken)
	userInfo, err := mysql.GetUserInfoByToken(c, nil, userToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetRoute GetUserInfoByToken failed, userInfo: %v", userInfo)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if userInfo.Status == helper.TransferredStatus {
		helper.BizResponse(c, http.StatusOK, helper.CodeInvalidUser, nil)
		return
	}

	routeInfo, err := mysql.GetRouteByToken(c, nil, routeToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetRoute GetRouteByToken failed, routeToken: %v", routeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if routeInfo.Status == helper.RouteDeleteStatus {
		helper.FormatLogPrint(helper.LOG, "GetRoute Route has deleted, routeToken: %v", routeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeMountDeleted, nil)
		return
	}

	// route鉴权
	if routeInfo.OwnerId != userInfo.ID {
		helper.BizResponse(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}

	// update status
	if routeInfo.Status != helper.RouteOpenStatus {
		nowTime := time.Now().Format("2006-01-02 15:04:05")
		err = mysql.UpdateRouteByToken(c, nil, routeToken, map[string]interface{}{
			"status":      helper.RouteOpenStatus,
			"update_time": nowTime,
		})
	}
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetRoute UpdateRouteByToken failed, err: %v", err)
	}

	routePointListOut, err := transform.CreateFmtRoutePointList(c, routeInfo.PointsList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetRoute CreateFmtRoutePointList failed, pointList: %v", routeInfo.PointsList)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	routeInfoOut := &out.RouteInfoOut{
		RouteToken:     routeInfo.RouteToken,
		RouteName:      routeInfo.RouteName,
		Status:         helper.RouteOpenStatus,
		RoutePointList: routePointListOut,
	}

	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, routeInfoOut)
}

func CloseRoute(c *gin.Context) {
	var closeRouteForm CloseRouteForm
	if err := c.ShouldBindJSON(&closeRouteForm); err != nil || closeRouteForm.RouteToken == "" {
		helper.FormatLogPrint(helper.WARNING, "CloseRoute bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "CloseRoute from: %+v", closeRouteForm)

	routeToken := closeRouteForm.RouteToken
	userToken := c.GetString(helper.UserToken)
	userInfo, err := mysql.GetUserInfoByToken(c, nil, userToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CloseRoute GetUserInfoByToken failed, userInfo: %v", userInfo)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if userInfo.Status == helper.TransferredStatus {
		helper.BizResponse(c, http.StatusOK, helper.CodeInvalidUser, nil)
		return
	}

	routeInfo, err := mysql.GetRouteByToken(c, nil, routeToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CloseRoute GetRouteByToken failed, routeToken: %v", routeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	if routeInfo.Status == helper.RouteDeleteStatus {
		helper.FormatLogPrint(helper.LOG, "CloseRoute Route has deleted, routeToken: %v", routeToken)
		helper.BizResponse(c, http.StatusOK, helper.CodeMountDeleted, nil)
		return
	}

	// route鉴权
	if routeInfo.OwnerId != userInfo.ID {
		helper.BizResponse(c, http.StatusOK, helper.CodeForbidden, nil)
		return
	}

	// update status
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	err = mysql.UpdateRouteByToken(c, nil, routeToken, map[string]interface{}{
		"status":      helper.RouteNormalStatus,
		"update_time": nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CloseRoute UpdateRouteByToken failed, err: %v", err)
	}

	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, nil)
}
