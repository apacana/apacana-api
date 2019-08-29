package transform

import (
	"encoding/json"
	"errors"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
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

func CreateFmtStrokeList(c *gin.Context, strokeStr string) (map[string]interface{}, error) {
	strokeList, err := StringToStrokeList(strokeStr)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "CreateFmtStrokeList StringToStrokeList failed, strokes: %v", strokeStr)
		return nil, errors.New("StringToStrokeList failed")
	}

	strokeInfos, err := mysql.MGetStrokeByID(c, nil, strokeList.HistoryStrokeList)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "GetUserInfo MGetStrokeByID failed, err: %v", err)
		return nil, errors.New("MGetStrokeByID failed")
	}

	defaultStroke := map[string]interface{}{}
	if strokeList.DefaultStroke != 0 {
		stroke, err := mysql.MGetStrokeByID(c, nil, []int64{strokeList.DefaultStroke})
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "GetUserInfo MGetStrokeByID failed, err: %v", err)
			return nil, errors.New("MGetStrokeByID failed")
		}
		if len(stroke) == 0 {
			helper.FormatLogPrint(helper.WARNING, "Stroke not exist, strokeID: %d", strokeList.DefaultStroke)
			return nil, errors.New("MGetStrokeByID failed")
		}
		defaultStroke = map[string]interface{}{
			"stroke_token": stroke[0].StrokeToken,
			"stroke_name":  stroke[0].StrokeName,
			"point_list":   stroke[0].PointsList,
			"update_time":  stroke[0].UpdateTime,
		}
	}

	historyStrokeInfoList := make([]map[string]interface{}, 0)
	for _, strokeInfo := range strokeInfos {
		if strokeInfo.Status == helper.StrokeDeleteStatus {
			continue
		}
		historyStrokeInfoList = append(historyStrokeInfoList, map[string]interface{}{
			"stroke_token": strokeInfo.StrokeToken,
			"stroke_name":  strokeInfo.StrokeName,
		})
	}
	strokeInfo := map[string]interface{}{
		"default_stroke":      defaultStroke,
		"history_stroke_list": historyStrokeInfoList,
	}
	return strokeInfo, nil
}
