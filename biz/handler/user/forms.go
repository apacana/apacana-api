package user

type RegisterUserForm struct {
	UserName string `json:"user_name" binding:"required"`
	PassWord string `json:"pass_word" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type LoginUserForm struct {
	UserName string `json:"user_name" binding:"required"`
	PassWord string `json:"pass_word" binding:"required"`
}

type HeartUserForm struct {
	Latitude  *float64 `json:"latitude" binding:"omitempty"`
	Longitude *float64 `json:"longitude" binding:"omitempty"`
	Zoom      *int64   `json:"zoom" binding:"omitempty"`
}
