package helper

import "errors"

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
	CodeSuccess          = int32(0)
	CodeFailed           = int32(1)
	CodeParmErr          = int32(2)
	CodeStrokeOutOfLimit = int32(3)
	CodeForbidden        = int32(4)
	CodeInvalidUser      = int32(5)
)

const (
	TouristStatus     = uint8(0)
	LoginUserStatus   = uint8(1)
	TransferredStatus = uint8(9)
)

const (
	StrokeNormalStatus = uint8(0)
	StrokeDeleteStatus = uint8(1)
)

var (
	ErrStrokeOutOfLimit = errors.New("ErrStrokeOutOfLimit")
	ErrInvaildCookie    = errors.New("ErrInvaildCookie")
)
