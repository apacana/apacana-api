package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RouteDirection struct {
	ID         int64  `gorm:"id" json:"id"`
	Direction  string `gorm:"direction" json:"direction"`
	RouteID    int64  `gorm:"stroke_id" json:"route_id"`
	Status     uint8  `gorm:"status" json:"status"`
	CreateTime string `gorm:"create_time" json:"create_time"`
	UpdateTime string `gorm:"update_time" json:"update_time"`
}

const (
	RouteDirectionTableName = "route_direction"
)

func (a *RouteDirection) TableName() string {
	return RouteDirectionTableName
}

func MGetRouteDirectionByID(c *gin.Context, tx *gorm.DB, IDs []int64) ([]*RouteDirection, error) {
	if tx == nil {
		tx = DB
	}
	var ref []*RouteDirection
	r := tx.Model(&RouteDirection{}).Where("id in (?)", IDs).Find(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func InsertRouteDirection(c *gin.Context, tx *gorm.DB, routeDirection *RouteDirection) error {
	return Insert(tx, routeDirection)
}
