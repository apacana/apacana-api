package helper

import "testing"

func Test_GenerateToken(t *testing.T) {
	for i := 0; i < 100; i += 1 {
		token := GenerateToken([]byte{'u', 's', 'e', 'r'}, "123456789")
		t.Log(token)
	}
}

func Test_SetCookie(t *testing.T) {
	token := GenerateToken([]byte{'u', 's', 'e', 'r'}, "123456789")
	cookie := SetCookie(token, SessionSalt)
	t.Log(cookie)
	cookie = SetCookie(token, SessionSalt)
	t.Log(cookie)
}
