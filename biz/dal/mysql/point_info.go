package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type PointInfo struct {
	ID         int64     `gorm:"id" json:"id"`
	PointID    string    `gorm:"point_id" json:"point_id"`
	PointType  PointType `gorm:"point_type" json:"point_type"`
	PointToken string    `gorm:"point_token" json:"point_token"`
	Text       string    `gorm:"text" json:"text"`
	PlaceName  string    `gorm:"place_name" json:"place_name"`
	Comment    string    `gorm:"comment" json:"comment"`
	Center     string    `gorm:"center" json:"center"`
	IconType   string    `gorm:"icon_type" json:"icon_type"`
	IconColor  string    `gorm:"icon_color" json:"icon_color"`
	Ext        string    `gorm:"ext" json:"ext"`
	Status     uint8     `gorm:"status" json:"status"`
	StrokeID   int64     `gorm:"stroke_id" json:"stroke_id"`
	CreateTime string    `gorm:"create_time" json:"create_time"`
	UpdateTime string    `gorm:"update_time" json:"update_time"`
}

type PointType uint8

const (
	PointType_UNKNOW      PointType = 0
	PointType_SEARCH      PointType = 1
	PointType_AGODA_HOTEL PointType = 2
)

const (
	PointInfoTableName = "point_info"
)

func (a *PointInfo) TableName() string {
	return PointInfoTableName
}

func MGetPointByID(c *gin.Context, tx *gorm.DB, IDs []int64) ([]*PointInfo, error) {
	if tx == nil {
		tx = DB
	}
	var ref []*PointInfo
	r := tx.Model(&PointInfo{}).Where("id in (?)", IDs).Find(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func GetPointByToken(c *gin.Context, tx *gorm.DB, pointToken string) (*PointInfo, error) {
	if tx == nil {
		tx = DB
	}
	var ref = &PointInfo{}
	r := tx.Model(&PointInfo{}).Where("point_token = ?", pointToken).First(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func AllowAddPoint(c *gin.Context, tx *gorm.DB, pointID string, pointType PointType, strokeID int64) (bool, error) {
	if tx == nil {
		tx = DB
	}
	var ref []*PointInfo
	r := tx.Model(&PointInfo{}).Where("stroke_id = ?", strokeID).Find(&ref)
	if r.Error != nil {
		if r.Error == gorm.ErrRecordNotFound {
			return true, nil
		}
		return false, r.Error
	}
	for _, point := range ref {
		if point.PointID == pointID && point.PointType == pointType {
			return false, nil
		}
	}
	return true, r.Error
}

func InsertPointInfo(c *gin.Context, tx *gorm.DB, pointInfo *PointInfo) error {
	return Insert(tx, pointInfo)
}
