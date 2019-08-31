package helper

import (
	"fmt"
	"github.com/apacana/apacana-api/biz/dal/mysql"
)

var pointTypeByName = map[string]mysql.PointType{
	"search_point": mysql.PointType_SEARCH,
	"agoda_hotel":  mysql.PointType_AGODA_HOTEL,
}

var nameByPointType = map[mysql.PointType]string{}

func init() {
	for mp, n := range pointTypeByName {
		nameByPointType[n] = mp
	}
}

func GetPointTypeByName(name string) (mysql.PointType, error) {
	if value, ok := pointTypeByName[name]; !ok {
		return mysql.PointType_UNKNOW, fmt.Errorf("unknown point type name: %s", name)
	} else {
		return value, nil
	}
}

func GetNameByPointType(pointType mysql.PointType) (string, error) {
	if value, ok := nameByPointType[pointType]; !ok {
		return "", fmt.Errorf("unknown point type: %d", pointType)
	} else {
		return value, nil
	}
}
