package route

type CreateRouteForm struct {
	RouteName   *string `json:"route_name" binding:"omitempty"`
	StrokeToken string  `json:"stroke_token" binding:"required"`
}

type AddRoutePointForm struct {
	RouteToken string  `json:"route_token" binding:"required"`
	PointToken string  `json:"point_token" binding:"required"`
	Direction  *string `json:"direction" binding:"omitempty"`
}

type CloseRouteForm struct {
	RouteToken string `json:"route_token" binding:"required"`
}

type OpenRouteForm struct {
	RouteToken string `json:"route_token" binding:"required"`
}
