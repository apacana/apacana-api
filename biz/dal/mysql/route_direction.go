package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RouteDirection struct {
	ID             int64         `gorm:"id" json:"id"`
	DirectionToken string        `gorm:"direction_token" json:"direction_token"`
	DirectionType  DirectionType `gorm:"direction_type" json:"direction_type"`
	Direction      string        `gorm:"direction" json:"direction"`
	RouteID        int64         `gorm:"stroke_id" json:"route_id"`
	Version        string        `gorm:"version" json:"version"`
	Status         uint8         `gorm:"status" json:"status"`
	CreateTime     string        `gorm:"create_time" json:"create_time"`
	UpdateTime     string        `gorm:"update_time" json:"update_time"`
}

const (
	RouteDirectionTableName = "route_direction"
)

type DirectionType uint8

const (
	DirectionType_DRIVINGTRAFFIC DirectionType = 0
	DirectionType_DRIVING        DirectionType = 1
	DirectionType_WALKING        DirectionType = 2
	DirectionType_CYCLING        DirectionType = 3
)

func (a *RouteDirection) TableName() string {
	return RouteDirectionTableName
}

func MGetRouteDirectionByID(c *gin.Context, tx *gorm.DB, IDs []int64) (map[int64]*RouteDirection, error) {
	if tx == nil {
		tx = DB
	}
	var ref []*RouteDirection
	r := tx.Model(&RouteDirection{}).Where("id in (?)", IDs).Find(&ref)
	if r.Error != nil {
		return nil, r.Error
	}

	refMap := make(map[int64]*RouteDirection, 0)
	for _, direction := range ref {
		refMap[direction.ID] = direction
	}
	return refMap, nil
}

func GetDirectionByToken(c *gin.Context, tx *gorm.DB, directionToken string) (*RouteDirection, error) {
	if tx == nil {
		tx = DB
	}
	var ref = &RouteDirection{}
	r := tx.Model(&RouteDirection{}).Where("direction_token = ?", directionToken).First(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func UpdateDirectionByID(c *gin.Context, tx *gorm.DB, id int64, attrs map[string]interface{}) error {
	if tx == nil {
		tx = DB
	}
	r := tx.Model(&RouteDirection{}).Where("id = ?", id).Update(attrs)
	return r.Error
}

func InsertRouteDirection(c *gin.Context, tx *gorm.DB, routeDirection *RouteDirection) (*RouteDirection, error) {
	err := Insert(tx, routeDirection)
	return routeDirection, err
}
