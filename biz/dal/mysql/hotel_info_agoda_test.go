package mysql

import "testing"

func Test_MGetHotelInfoAgodaByHotelID(t *testing.T) {
	InitMysql()
	resp, err := MGetHotelInfoAgodaByHotelID(nil, nil, []int64{1})
	t.Logf("%+v, %v", *resp[0], err)
}
