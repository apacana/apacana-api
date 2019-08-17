package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type StrokeInfo struct {
	ID          int64  `gorm:"id" json:"id"`
	StrokeToken string `gorm:"stroke_token" json:"stroke_token"`
	StrokeName  string `gorm:"stroke_name" json:"stroke_name"`
	RoutesList  string `gorm:"routes_list" json:"routes_list"`
	OwnerID     int64  `gorm:"owner_id" json:"owner_id"`
	Status      uint8  `gorm:"status" json:"status"`
	CreateTime  string `gorm:"create_time" json:"create_time"`
	UpdateTime  string `gorm:"update_time" json:"update_time"`
}

const (
	StrokeInfoTableName = "stroke_info"
)

func (a *StrokeInfo) TableName() string {
	return StrokeInfoTableName
}

func MGetStrokeByID(c *gin.Context, tx *gorm.DB, IDs []int64) ([]*StrokeInfo, error) {
	if tx == nil {
		tx = DB
	}
	var ref []*StrokeInfo
	r := tx.Model(&StrokeInfo{}).Where("id in (?)", IDs).Find(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func GetStrokeByToken(c *gin.Context, tx *gorm.DB, strokeToken string) (*StrokeInfo, error) {
	if tx == nil {
		tx = DB
	}
	var ref = &StrokeInfo{}
	r := tx.Model(&StrokeInfo{}).Where("stroke_token = ?", strokeToken).First(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func UpdateStrokeByToken(c *gin.Context, tx *gorm.DB, strokeToken string, attrs map[string]interface{}) error {
	if tx == nil {
		tx = DB
	}
	r := tx.Model(&StrokeInfo{}).Where("stroke_token = ?", strokeToken).Update(attrs)
	return r.Error
}

func InsertStrokeInfo(c *gin.Context, tx *gorm.DB, strokeInfo *StrokeInfo) error {
	return Insert(tx, strokeInfo)
}

func ChangeStrokeOwner(c *gin.Context, tx *gorm.DB, originalID int64, purposeID int64) error {
	if tx == nil {
		tx = DB
	}
	r := tx.Model(&StrokeInfo{}).Where("owner_id = ?", originalID).Update(map[string]interface{}{
		"owner_id": purposeID,
	})
	return r.Error
}
