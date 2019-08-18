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

func StringToStrokeList(str string) (*StrokeList, error) {
	strokeList := &StrokeList{
		StrokeList: make([]int64, 0),
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

func CreateFmtStrokeList(c *gin.Context, strokeStr string) ([]map[string]interface{}, error) {
	strokeList, err := StringToStrokeList(strokeStr)
	if err != nil {
		FormatLogPrint(ERROR, "GetUserInfo StringToStrokeList failed, strokes: %v", strokeStr)
		return nil, errors.New("StringToStrokeList failed")
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
