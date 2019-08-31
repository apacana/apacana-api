package transform

import (
	"encoding/json"
	"errors"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"unsafe"
)

type RouteList struct {
	RouteList []int64 `json:"route_list"`
	Ext       string  `json:"ext"`
}

func StringToRouteList(str string) (*RouteList, error) {
	routeList := &RouteList{
		RouteList: make([]int64, 0),
	}
	if len(str) != 0 {
		var newRouteList RouteList
		err := json.Unmarshal([]byte(str), &newRouteList)
		if err != nil {
			return nil, errors.New("StringToRouteList failed")
		}
		routeList = &newRouteList
	}
	return routeList, nil
}

func PackRouteList(routeList *RouteList) *string {
	bytesData, _ := json.Marshal(*routeList)
	return (*string)(unsafe.Pointer(&bytesData))
}

func CreateFmtRouteList(c *gin.Context, routesStr string) ([]map[string]interface{}, error) {
	routeList, err := StringToRouteList(routesStr)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtRouteList StringToRouteList failed, routes: %v", routesStr)
		return nil, errors.New("CreateFmtRouteList failed")
	}

	RouteInfos, err := mysql.MGetRouteByID(c, nil, routeList.RouteList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetUserInfo MGetStrokeByID failed, err: %v", err)
		return nil, errors.New("CreateFmtRouteList failed")
	}
	routeInfoList := make([]map[string]interface{}, 0)
	for _, routeInfo := range RouteInfos {
		routeInfoList = append(routeInfoList, map[string]interface{}{
			"route_token": routeInfo.RouteToken,
			"route_name":  routeInfo.RouteName,
		})
	}
	return routeInfoList, nil
}
