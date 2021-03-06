package stroke

type CreateStrokeForm struct {
	StrokeName *string `json:"stroke_name" binding:"omitempty"`
}

type ChangeDefaultForm struct {
	StrokeToken string `json:"stroke_token" binding:"required"`
}

type UpdateStrokeForm struct {
	StrokeToken string  `json:"stroke_token" binding:"required"`
	StrokeName  *string `json:"stroke_name" binding:"omitempty"`
}
