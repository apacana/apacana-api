package out

/*User*/

type UserInfoOut struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	Status uint8  `json:"status"`
}

/*Stroke*/

type DefaultStrokeOut struct {
	StrokeToken string          `json:"stroke_token"`
	StrokeName  string          `json:"stroke_name"`
	PointList   []*PointInfoOut `json:"point_list"`
	RouteList   []*RouteInfoOut `json:"route_list"`
	UpdateTime  string          `json:"update_time"`
}

type StrokeUpdateOut struct {
	StrokeToken string `json:"stroke_token"`
	StrokeName  string `json:"stroke_name"`
	UpdateTime  string `json:"update_time"`
}

type StrokeInfoOut struct {
	StrokeToken string `json:"stroke_token"`
	StrokeName  string `json:"stroke_name"`
}

type StrokesInfoOut struct {
	DefaultStroke     *DefaultStrokeOut `json:"default_stroke"`
	HistoryStrokeList []*StrokeInfoOut  `json:"history_stroke_list"`
}

/*Route*/

type RouteInfoOut struct {
	RouteToken     string        `json:"route_token"`
	RouteName      string        `json:"route_name"`
	Status         uint8         `json:"status"`
	RoutePointList []*RoutePoint `json:"route_point"`
	UpdateTime     string        `json:"update_time"`
}

type RoutePoint struct {
	PointID   string `json:"point_id"`
	PointType string `json:"point_type"`
	Text      string `json:"text"`
	Direction string `json:"direction"`
}

/*Point*/

type PointInfoOut struct {
	PointToken string `json:"point_token"`
	PointID    string `json:"point_id"`
	PointType  string `json:"point_type"`
	Text       string `json:"text"`
	PlaceName  string `json:"place_name"`
	Center     string `json:"center"`
	Comment    string `json:"comment"`
	IconType   string `json:"icon_type"`
	IconColor  string `json:"icon_color"`
	Ext        string `json:"ext"`
}
