package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UserInfo struct {
	ID         int64  `gorm:"id" json:"id"`
	Token      string `gorm:"token" json:"token"`
	UserName   string `gorm:"user_name" json:"user_name"`
	PassWord   string `gorm:"pass_word" json:"pass_word"`
	Name       string `gorm:"name" json:"name"`
	Strokes    string `gorm:"strokes" json:"strokes"`
	Status     uint8  `gorm:"status" json:"status"`
	CreateTime string `gorm:"create_time" json:"create_time"`
	UpdateTime string `gorm:"update_time" json:"update_time"`
}

const (
	UserInfoTableName = "user_info"
)

func (a *UserInfo) TableName() string {
	return UserInfoTableName
}

func GetUserInfoByToken(c *gin.Context, tx *gorm.DB, token string) (*UserInfo, error) {
	if tx == nil {
		tx = DB
	}
	var ref = &UserInfo{}
	r := tx.Model(&UserInfo{}).Where("token = ?", token).First(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func GetUserByUserPassWord(c *gin.Context, tx *gorm.DB, userName string, passWord string) (*UserInfo, error) {
	if tx == nil {
		tx = DB
	}
	var ref = &UserInfo{}
	r := tx.Model(&UserInfo{}).Where("user_name = ? AND pass_word = ?", userName, passWord).First(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func UpdateUserInfo(c *gin.Context, tx *gorm.DB, id int64, attrs map[string]interface{}) error {
	if tx == nil {
		tx = DB
	}
	r := tx.Model(&UserInfo{}).Where("id = ?", id).Update(attrs)
	return r.Error
}

func InsertUserInfo(c *gin.Context, userInfo *UserInfo) error {
	return Insert(nil, userInfo)
}
