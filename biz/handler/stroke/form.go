package stroke

type CreateStrokeForm struct {
	StrokeName string `json:"stroke_name" binding:"required"`
}
