package helper

import (
	"encoding/json"
	"errors"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/gin-gonic/gin"
	"unsafe"
)

type StrokeList struct {
	StrokeList []int64 `json:"stroke_list"`
	Ext        string  `json:"ext"`
}

func StringToStrokeList(str string) *StrokeList {
	// todo: 代码优化
	var strikeList StrokeList
	err := json.Unmarshal([]byte(str), &strikeList)
	if err != nil {
		FormatLogPrint(ERROR, "StringToStrokeList failed, err: %v, str: %v", err, str)
		return nil
	}
	return &strikeList
}

func PackStrokeList(strikeList *StrokeList) *string {
	bytesData, _ := json.Marshal(*strikeList)
	return (*string)(unsafe.Pointer(&bytesData))
}

func CreateFmtStrokeList(c *gin.Context, strokeStr string) ([]map[string]interface{}, error) {
	strokeList := &StrokeList{
		StrokeList: make([]int64, 0),
	}
	if len(strokeStr) != 0 {
		strokeList = StringToStrokeList(strokeStr)
		if strokeList == nil {
			FormatLogPrint(ERROR, "GetUserInfo StringToStrokeList failed, strokes: %v", strokeStr)
			return nil, errors.New("StringToStrokeList failed")
		}
	}

	strokeInfos, err := mysql.MGetStrokeByID(c, nil, strokeList.StrokeList)
	if err != nil {
		FormatLogPrint(ERROR, "GetUserInfo MGetStrokeByID failed, err: %v", err)
		return nil, errors.New("MGetStrokeByID failed")
	}
	strokeInfoList := make([]map[string]interface{}, 0)
	for _, strokeInfo := range strokeInfos {
		strokeInfoList = append(strokeInfoList, map[string]interface{}{
			"stroke_token": strokeInfo.StrokeToken,
			"stroke_name":  strokeInfo.StrokeName,
		})
	}
	return strokeInfoList, nil
}
