package helper

const YearTime = 31536000
const SessionSalt = "saltjKNlvO8ybWIqkuZCCi1jgd"

const (
	ReportHttpCode = "REPORT_HTTP_CODE"
	ReportBizCode  = "REPORT_BIZ_CODE"
	UserToken      = "USER_TOKEN"
	ApacanaSession = "apacana-session"
)

const (
	LOG     = "[LOG]"
	WARNING = "[WARNING]"
	ERROR   = "[ERROR]"
)

const (
	CodeSuccess   = int32(0)
	CodeFailed    = int32(1)
	CodeForbidden = int32(4)
)

const (
	TouristStatus   = uint8(0)
	LoginUserStatus = uint8(1)
)
