package route

import (
	"github.com/apacana/apacana-api/biz/config"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/out"
	"github.com/apacana/apacana-api/biz/transform"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRoute(c *gin.Context) {
	var createRouteForm CreateRouteForm
	if err := c.ShouldBindJSON(&createRouteForm); err != nil ||
		len(createRouteForm.RouteName) == 0 || len(createRouteForm.RouteName) > 24 {
		helper.FormatLogPrint(helper.WARNING, "CreateStroke bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "CreateRoute from: %+v", createRouteForm)

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
	routeToken, err := createStrokeRoute(c, strokeInfo, routeList, createRouteForm.RouteName)
	if err != nil {
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, map[string]interface{}{
		"route_name":  createRouteForm.RouteName,
		"route_token": routeToken,
	})
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, &out.RouteInfoOut{
		RouteToken: routeToken,
		RouteName:  createRouteForm.RouteName,
	})
}
