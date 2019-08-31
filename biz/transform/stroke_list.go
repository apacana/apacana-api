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

type StrokeList struct {
	DefaultStroke     int64   `json:"default_stroke"`
	HistoryStrokeList []int64 `json:"history_stroke_list"`
	Ext               string  `json:"ext"`
}

func StringToStrokeList(str string) (*StrokeList, error) {
	strokeList := &StrokeList{
		DefaultStroke:     0,
		HistoryStrokeList: make([]int64, 0),
	}
	if len(str) != 0 {
		var newStrikeList StrokeList
		err := json.Unmarshal([]byte(str), &newStrikeList)
		if err != nil {
			return nil, errors.New("StringToStrokeList failed")
		}
		strokeList = &newStrikeList
	}
	return strokeList, nil
}

func PackStrokeList(strikeList *StrokeList) *string {
	bytesData, _ := json.Marshal(*strikeList)
	return (*string)(unsafe.Pointer(&bytesData))
}

func CreateFmtStrokeList(c *gin.Context, strokeStr string) (*out.StrokesInfoOut, error) {
	strokesInfoOut := &out.StrokesInfoOut{
		HistoryStrokeList: make([]*out.StrokeInfoOut, 0),
	}
	strokeList, err := StringToStrokeList(strokeStr)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtStrokeList StringToStrokeList failed, strokes: %v", strokeStr)
		return nil, errors.New("StringToStrokeList failed")
	}

	strokeInfos, err := mysql.MGetStrokeByID(c, nil, strokeList.HistoryStrokeList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtStrokeList MGetStrokeByID failed, err: %v", err)
		return nil, errors.New("MGetStrokeByID failed")
	}

	if strokeList.DefaultStroke != 0 {
		stroke, err := mysql.MGetStrokeByID(c, nil, []int64{strokeList.DefaultStroke})
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "CreateFmtStrokeList MGetStrokeByID failed, err: %v", err)
			return nil, errors.New("MGetStrokeByID failed")
		}
		if len(stroke) == 0 {
			helper.FormatLogPrint(helper.WARNING, "Stroke not exist, strokeID: %d", strokeList.DefaultStroke)
			return nil, errors.New("MGetStrokeByID failed")
		}
		routeListOut, err := CreateFmtRouteList(c, stroke[0].RoutesList)
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "CreateFmtStrokeList CreateFmtRouteList failed, err: %v", err)
			return nil, errors.New("CreateFmtRouteList failed")
		}
		pointListOut, err := CreateFmtPointList(c, stroke[0].PointsList)
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "CreateFmtStrokeList CreateFmtPointList failed, err: %v", err)
			return nil, errors.New("CreateFmtPointList failed")
		}
		strokesInfoOut.DefaultStroke = &out.DefaultStrokeOut{
			StrokeToken: stroke[0].StrokeToken,
			StrokeName:  stroke[0].StrokeName,
			PointList:   pointListOut,
			RouteList:   routeListOut,
			UpdateTime:  stroke[0].UpdateTime,
		}
	}

	for _, strokeInfo := range strokeInfos {
		strokesInfoOut.HistoryStrokeList = append(strokesInfoOut.HistoryStrokeList, &out.StrokeInfoOut{
			StrokeToken: strokeInfo.StrokeToken,
			StrokeName:  strokeInfo.StrokeName,
		})
	}
	return strokesInfoOut, nil
}
