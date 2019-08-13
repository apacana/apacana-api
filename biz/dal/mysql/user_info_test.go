package mysql

import "testing"

func Test_GetUserInfoByToken(t *testing.T) {
	InitMysql()
	resp, err := GetUserInfoByToken(nil, nil, "test")
	t.Logf("%+v, %v", *resp, err)
	resp, err = GetUserInfoByToken(nil, nil, "nibaba")
	t.Logf("%+v, %v", resp, err)
}
