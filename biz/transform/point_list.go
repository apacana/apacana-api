package transform

import (
	"encoding/json"
	"errors"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"unsafe"
)

type PointList struct {
	PointList []int64 `json:"point_list"`
	Ext       string  `json:"ext"`
}

func StringToPointList(str string) (*PointList, error) {
	pointRoute := &PointList{
		PointList: make([]int64, 0),
	}
	if len(str) != 0 {
		var newPointRoute PointList
		err := json.Unmarshal([]byte(str), &newPointRoute)
		if err != nil {
			return nil, errors.New("StringToPointList failed")
		}
		pointRoute = &newPointRoute
	}
	return pointRoute, nil
}

func PackPointList(pointRoute *PointList) *string {
	bytesData, _ := json.Marshal(*pointRoute)
	return (*string)(unsafe.Pointer(&bytesData))
}

func CreateFmtPointList(c *gin.Context, pointsStr string) ([]map[string]interface{}, error) {
	pointList, err := StringToPointList(pointsStr)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtPointList StringToPointList failed, points: %v", pointsStr)
		return nil, errors.New("CreateFmtPointList failed")
	}

	pointInfos, err := mysql.MGetPointByID(c, nil, pointList.PointList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtPointList MGetPointByID failed, err: %v", err)
		return nil, errors.New("CreateFmtPointList failed")
	}
	pointInfoList := make([]map[string]interface{}, 0)
	for _, pointInfo := range pointInfos {
		pointTypeName, _ := helper.GetNameByPointType(pointInfo.PointType)
		pointInfoList = append(pointInfoList, map[string]interface{}{
			"point_id":    pointInfo.PointID,
			"point_type":  pointTypeName,
			"point_token": pointInfo.PointToken,
			"text":        pointInfo.Text,
			"place_name":  pointInfo.PlaceName,
			"center":      pointInfo.Center,
			"comment":     pointInfo.Comment,
			"icon_type":   pointInfo.IconType,
			"icon_color":  pointInfo.IconColor,
			"status":      pointInfo.Status,
			"ext":         pointInfo.Ext,
		})
	}
	return pointInfoList, nil
}
