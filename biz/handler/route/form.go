package route

type CreateRouteForm struct {
	RouteName   *string `json:"route_name" binding:"omitempty"`
	StrokeToken string  `json:"stroke_token" binding:"required"`
}

type AddRoutePointForm struct {
	RouteToken    string  `json:"route_token" binding:"required"`
	PointToken    string  `json:"point_token" binding:"required"`
	DirectionType *string `json:"direction_type" binding:"omitempty"`
	Direction     *string `json:"direction" binding:"omitempty"`
}

type RemoveRoutePointForm struct {
	RouteToken    string  `json:"route_token" binding:"required"`
	Index         *int    `json:"index" binding:"exists"`
	DirectionType *string `json:"direction_type" binding:"omitempty"`
	Direction     *string `json:"direction" binding:"omitempty"`
}

type CloseRouteForm struct {
	RouteToken string `json:"route_token" binding:"required"`
}

type OpenRouteForm struct {
	RouteToken string `json:"route_token" binding:"required"`
}

type UpdateDirectionForm struct {
	Index         int     `json:"index" binding:"required"`
	RouteToken    string  `json:"route_token" binding:"required"`
	DirectionType *string `json:"direction_type" binding:"omitempty"`
	Direction     *string `json:"direction" binding:"omitempty"`
}

type UpdateRouteForm struct {
	RouteToken string  `json:"route_token" binding:"required"`
	RouteName  *string `json:"route_name" binding:"omitempty"`
}
