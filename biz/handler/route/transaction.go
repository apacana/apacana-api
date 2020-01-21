package route

import (
	"errors"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/transform"
	"github.com/gin-gonic/gin"
	"time"
)

func createStrokeRoute(c *gin.Context, strokeInfo *mysql.StrokeInfo, routeList *transform.RouteList, routeName string, routeColor string) (routeToken string, nowTime string, err error) {
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
	routeInfo, err := mysql.InsertRouteInfo(c, tx, &mysql.RouteInfo{
		RouteToken: routeToken,
		RouteName:  routeName,
		RouteColor: routeColor,
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

	// update stroke
	routeList.RouteList = append(routeList.RouteList, routeInfo.ID)
	err = mysql.UpdateStrokeByToken(c, tx, strokeInfo.StrokeToken, map[string]interface{}{
		"routes_list": *transform.PackRouteList(routeList),
	})

	return
}

func addRoutePoint(c *gin.Context, routeInfo *mysql.RouteInfo, pointInfo *mysql.PointInfo, addRoutePointForm AddRoutePointForm) (nowTime string, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "addRoutePoint failed, err: %v", err)
			}
		}
	}()

	// update route
	routePointList, err := transform.StringToRoutePointList(routeInfo.PointsList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "addRoutePoint StringToRoutePointList failed, err: %v", err)
		return
	}
	routePointList.PointList = append(routePointList.PointList, pointInfo.ID)

	// update direction
	var directionID int64 = 0
	nowTime = time.Now().Format("2006-01-02 15:04:05")
	if addRoutePointForm.Direction != nil && *addRoutePointForm.Direction != "" {
		directionType := mysql.DirectionType_DRIVINGTRAFFIC
		if addRoutePointForm.DirectionType != nil && *addRoutePointForm.DirectionType != "" {
			aDirectionType, err := helper.GetDirectionTypeByName(*addRoutePointForm.DirectionType)
			if err == nil {
				directionType = aDirectionType
			}
		}
		directToken := helper.GenerateToken([]byte{'d', 'i', 'r', 'e', 'c', 't'}, "")
		directionInfo, err := mysql.InsertRouteDirection(c, tx, &mysql.RouteDirection{
			DirectionToken: directToken,
			DirectionType:  directionType,
			Direction:      *addRoutePointForm.Direction,
			RouteID:        routeInfo.ID,
			Version:        "v1",
			CreateTime:     nowTime,
			UpdateTime:     nowTime,
		})
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "addRoutePoint InsertRouteDirection failed, err: %v", err)
		} else {
			directionID = directionInfo.ID
		}
	}
	routePointList.DirectionList = append(routePointList.DirectionList, directionID)

	err = mysql.UpdateRouteByToken(c, tx, routeInfo.RouteToken, map[string]interface{}{
		"update_time": nowTime,
		"points_list": *transform.PackRoutePointList(routePointList),
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "addRoutePoint UpdateRouteByToken failed, err: %v", err)
		return
	}
	err = mysql.UpdateStrokeByID(c, tx, routeInfo.StrokeID, map[string]interface{}{
		"update_time": nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "addRoutePoint UpdateStrokeByID failed, err: %v", err)
		return
	}

	return
}

func updateRoute(c *gin.Context, nowTime string, routeInfo *mysql.RouteInfo, attrs map[string]interface{}) (err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "addRoutePoint failed, err: %v", err)
			}
		}
	}()

	attrs["update_time"] = nowTime
	err = mysql.UpdateRouteByToken(c, tx, routeInfo.RouteToken, attrs)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "updateRoute UpdateRouteByToken failed, err: %v", err)
		return
	}
	err = mysql.UpdateStrokeByID(c, tx, routeInfo.StrokeID, map[string]interface{}{
		"update_time": nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "updateRoute UpdateStrokeByID failed, err: %v", err)
		return
	}

	return
}

func updateDirection(c *gin.Context, routeInfo *mysql.RouteInfo, updateDirectionForm UpdateDirectionForm) (nowTime string, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if err == nil {
			err = tx.Commit().Error
		}
		if err != nil {
			if r := tx.Rollback(); r.Error != nil {
				helper.FormatLogPrint(helper.ERROR, "updateDirection failed, err: %v", err)
			}
		}
	}()
	routePointList, err := transform.StringToRoutePointList(routeInfo.PointsList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "updateDirection StringToRoutePointList failed, err: %v", err)
		return
	}
	if updateDirectionForm.Index >= len(routePointList.DirectionList) {
		helper.FormatLogPrint(helper.WARNING, "updateDirection routePointList DirectionList out of index, index: %d", updateDirectionForm.Index)
		err = errors.New("updateDirection routePointList DirectionList out of index")
		return
	}

	nowTime = time.Now().Format("2006-01-02 15:04:05")
	directionType := mysql.DirectionType_DRIVINGTRAFFIC
	if updateDirectionForm.DirectionType != nil && *updateDirectionForm.DirectionType != "" {
		aDirectionType, err := helper.GetDirectionTypeByName(*updateDirectionForm.DirectionType)
		if err == nil {
			directionType = aDirectionType
		}
	}
	direction := ""
	if updateDirectionForm.Direction != nil && *updateDirectionForm.Direction != "" {
		direction = *updateDirectionForm.Direction
	}
	err = mysql.UpdateDirectionByID(c, nil, routePointList.DirectionList[updateDirectionForm.Index], map[string]interface{}{
		"direction_type": directionType,
		"direction":      direction,
		"update_time":    nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "updateDirection UpdateDirectionByID failed, err: %v", err)
		return
	}

	err = mysql.UpdateRouteByToken(c, nil, routeInfo.RouteToken, map[string]interface{}{
		"update_time": nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "updateDirection UpdateRouteByToken failed, err: %v", err)
		return
	}
	err = mysql.UpdateStrokeByID(c, nil, routeInfo.StrokeID, map[string]interface{}{
		"update_time": nowTime,
	})
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "updateDirection UpdateStrokeByID failed, err: %v", err)
		return
	}

	return
}
