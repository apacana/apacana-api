package user

type CreateTouristForm struct {
	Key string `json:"key" binding:"required"`
}

type RegisterUserForm struct {
	UserName string `json:"user_name" binding:"required"`
	PassWord string `json:"pass_word" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type LoginUserForm struct {
	UserName string `json:"user_name" binding:"required"`
	PassWord string `json:"pass_word" binding:"required"`
}
