package point

type AddPointForm struct {
	PointID   string  `json:"point_id" binding:"required"`
	PointType string  `json:"point_type" binding:"required"`
	Text      string  `json:"text" binding:"required"`
	PlaceName *string `json:"place_name" binding:"omitempty"`
	Center    string  `json:"center" binding:"required"`
	IconType  *string `json:"icon_type" binding:"omitempty"`
	IconColor *string `json:"icon_color" binding:"omitempty"`
	Ext       *string `json:"ext" binding:"omitempty"`
}

type DeletePointForm struct {
	PointID   string `json:"point_id" binding:"required"`
	PointType string `json:"point_type" binding:"required"`
}
