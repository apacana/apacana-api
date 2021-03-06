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
	Center     string `gorm:"center" json:"center"`
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
	r := tx.Model(&UserInfo{}).Where("user_name = ? AND pass_word = ? AND status = 1", userName, passWord).First(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func UserNameHasExist(c *gin.Context, tx *gorm.DB, userName string) bool {
	if tx == nil {
		tx = DB
	}
	var ref = &UserInfo{}
	r := tx.Model(&UserInfo{}).Where("user_name = ?", userName).First(&ref)
	if r.Error != nil && r.Error == gorm.ErrRecordNotFound {
		return false
	}
	return true
}

func UpdateUserInfo(c *gin.Context, tx *gorm.DB, id int64, attrs map[string]interface{}) error {
	if tx == nil {
		tx = DB
	}
	r := tx.Model(&UserInfo{}).Where("id = ?", id).Update(attrs)
	return r.Error
}

func UpdateUserInfoByToken(c *gin.Context, tx *gorm.DB, token string, attrs map[string]interface{}) error {
	if tx == nil {
		tx = DB
	}
	r := tx.Model(&UserInfo{}).Where("token = ?", token).Update(attrs)
	return r.Error
}

func InsertUserInfo(c *gin.Context, userInfo *UserInfo) (*UserInfo, error) {
	err := Insert(nil, userInfo)
	return userInfo, err
}
