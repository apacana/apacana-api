package route

type CreateRouteForm struct {
	RouteName   string `json:"route_name" binding:"required"`
	StrokeToken string `json:"stroke_token" binding:"required"`
}
