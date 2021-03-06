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
	CodeSuccess      = int32(0)
	CodeFailed       = int32(1)
	CodeParmErr      = int32(2)
	CodeForbidden    = int32(4)
	CodeMountDeleted = int32(6)

	CodeInvalidUser = int32(1001)

	CodeStrokeOutOfLimit = int32(2001)

	CodeRouteOutOfLimit = int32(3001)

	CodePointOutOfLimit = int32(4001)
	CodePointExist      = int32(4002)
	CodePointUsed       = int32(4003)
)

const (
	TouristStatus     = uint8(0)
	LoginUserStatus   = uint8(1)
	TransferredStatus = uint8(9)
)

const (
	StrokeNormalStatus = uint8(0)
	StrokeDeleteStatus = uint8(8)
)

const (
	RouteNormalStatus = uint8(0)
	RouteOpenStatus   = uint8(1)
	RouteDeleteStatus = uint8(8)
)

const (
	PointNormalStatus = uint8(0)
	PointDeleteStatus = uint8(8)
)

var (
	ErrStrokeOutOfLimit = errors.New("ErrStrokeOutOfLimit")
	ErrInvaildCookie    = errors.New("ErrInvaildCookie")
)
