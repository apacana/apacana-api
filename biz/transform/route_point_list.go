package transform

import (
	"encoding/json"
	"errors"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/apacana/apacana-api/biz/out"
	"github.com/gin-gonic/gin"
	"unsafe"
)

type RoutePointList struct {
	PointList     []int64 `json:"point_list"`
	DirectionList []int64 `json:"direction_list"`
	Ext           string  `json:"ext"`
}

func StringToRoutePointList(str string) (*RoutePointList, error) {
	routePointList := &RoutePointList{
		PointList:     make([]int64, 0),
		DirectionList: make([]int64, 0),
	}
	if len(str) != 0 {
		var newRoutePointList RoutePointList
		err := json.Unmarshal([]byte(str), &newRoutePointList)
		if err != nil {
			return nil, errors.New("StringToRoutePointList failed")
		}
		routePointList = &newRoutePointList
	}
	return routePointList, nil
}

func PackRoutePointList(routePointList *RoutePointList) *string {
	bytesData, _ := json.Marshal(*routePointList)
	return (*string)(unsafe.Pointer(&bytesData))
}

func CreateFmtRoutePointList(c *gin.Context, routePointStr string) ([]*out.RoutePoint, error) {
	routePointList, err := StringToRoutePointList(routePointStr)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtRoutePointList StringToRoutePointList failed, routePointStr: %v", routePointStr)
		return nil, errors.New("CreateFmtRoutePointList failed")
	}
	routePointOut := make([]*out.RoutePoint, len(routePointList.PointList))
	if len(routePointOut) == 0 {
		return routePointOut, nil
	}

	pointMap, _, err := mysql.MGetPointByID(c, nil, routePointList.PointList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtRoutePointList MGetPointByID failed, err: %v", err)
		return nil, errors.New("CreateFmtRoutePointList failed")
	}
	directionMap, err := mysql.MGetRouteDirectionByID(c, nil, routePointList.DirectionList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtRoutePointList MGetRouteDirectionByID failed, err: %v", err)
		return nil, errors.New("CreateFmtRoutePointList failed")
	}

	for index, id := range routePointList.PointList {
		point := pointMap[id]
		pointTypeName, _ := helper.GetNameByPointType(point.PointType)
		direction := ""
		directionType := "driving-traffic"
		if direct, ok := directionMap[routePointList.DirectionList[index]]; ok {
			direction = direct.Direction
			aDirectionType, err := helper.GetNameByDirectionType(direct.DirectionType)
			if err == nil {
				directionType = aDirectionType
			}
		}

		routePointOut[index] = &out.RoutePoint{
			PointID:       point.PointID,
			PointType:     pointTypeName,
			Text:          point.Text,
			Center:        point.Center,
			Direction:     direction,
			DirectionType: directionType,
		}
	}

	return routePointOut, nil
}
