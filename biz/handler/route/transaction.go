package route

import (
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/transform"
	"github.com/gin-gonic/gin"
	"time"
)

func createStrokeRoute(c *gin.Context, strokeInfo *mysql.StrokeInfo, routeList *transform.RouteList, routeName string) (routeToken string, nowTime string, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "createStrokeRoute failed, err: %v", err)
			}
		}
	}()

	// create route
	nowTime = time.Now().Format("2006-01-02 15:04:05")
	routeToken = helper.GenerateToken([]byte{'r', 'o', 'u', 't', 'e'}, "")
	err = mysql.InsertRouteInfo(c, tx, &mysql.RouteInfo{
		RouteToken: routeToken,
		RouteName:  routeName,
		StrokeID:   strokeInfo.ID,
		OwnerId:    strokeInfo.OwnerID,
		Status:     helper.RouteOpenStatus,
		CreateTime: nowTime,
		UpdateTime: nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "createStrokeRoute InsertRouteInfo failed, err: %v", err)
		return
	}
	routeInfo, err := mysql.GetRouteByToken(c, tx, routeToken)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "createStrokeRoute GetRouteByToken failed, err: %v, routeToken: %v", err, routeToken)
		return
	}

	// update stroke
	routeList.RouteList = append(routeList.RouteList, routeInfo.ID)
	err = mysql.UpdateStrokeByToken(c, tx, strokeInfo.StrokeToken, map[string]interface{}{
		"routes_list": *transform.PackRouteList(routeList),
	})

	return
}
