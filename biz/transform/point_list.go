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

func CreateFmtPointList(c *gin.Context, pointsStr string) ([]*out.PointInfoOut, error) {
	pointListOut := make([]*out.PointInfoOut, 0)
	pointList, err := StringToPointList(pointsStr)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtPointList StringToPointList failed, points: %v", pointsStr)
		return nil, errors.New("CreateFmtPointList failed")
	}

	if len(pointList.PointList) == 0 {
		return make([]*out.PointInfoOut, 0), nil
	}
	pointInfos, err := mysql.MGetPointByID(c, nil, pointList.PointList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtPointList MGetPointByID failed, err: %v", err)
		return nil, errors.New("CreateFmtPointList failed")
	}
	for _, pointInfo := range pointInfos {
		pointTypeName, _ := helper.GetNameByPointType(pointInfo.PointType)
		pointListOut = append(pointListOut, &out.PointInfoOut{
			PointID:    pointInfo.PointID,
			PointType:  pointTypeName,
			PointToken: pointInfo.PointToken,
			Text:       pointInfo.Text,
			PlaceName:  pointInfo.PlaceName,
			Center:     pointInfo.Center,
			Comment:    pointInfo.Comment,
			IconType:   pointInfo.IconType,
			IconColor:  pointInfo.IconColor,
			Ext:        pointInfo.Ext,
		})
	}
	return pointListOut, nil
}
