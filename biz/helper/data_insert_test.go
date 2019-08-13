package helper

import (
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"testing"
)

func Test_insertHotelInfoAgoda(t *testing.T) {
	mysql.InitMysql()
	insertHotelInfoAgoda("/Users/lupengyu/Desktop/data/file-info-agoda.csv")
}
