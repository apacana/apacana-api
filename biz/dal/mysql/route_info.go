package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type RouteInfo struct {
	ID         int64  `gorm:"id" json:"id"`
	RouteToken string `gorm:"route_token" json:"route_token"`
	RouteName  string `gorm:"route_name" json:"route_name"`
	RouteColor string `gorm:"route_color" json:"route_color"`
	PointsList string `gorm:"points_list" json:"points_list"`
	StrokeID   int64  `gorm:"stroke_id" json:"stroke_id"`
	OwnerId    int64  `gorm:"owner_id" json:"user_id"`
	Status     uint8  `gorm:"status" json:"status"`
	CreateTime string `gorm:"create_time" json:"create_time"`
	UpdateTime string `gorm:"update_time" json:"update_time"`
}

const (
	RouteInfoTableName = "route_info"
)

func (a *RouteInfo) TableName() string {
	return RouteInfoTableName
}

func MGetRouteByID(c *gin.Context, tx *gorm.DB, IDs []int64) ([]*RouteInfo, error) {
	if tx == nil {
		tx = DB
	}
	var ref []*RouteInfo
	r := tx.Model(&RouteInfo{}).Where("id in (?)", IDs).Find(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func GetRouteByToken(c *gin.Context, tx *gorm.DB, routeToken string) (*RouteInfo, error) {
	if tx == nil {
		tx = DB
	}
	var ref = &RouteInfo{}
	r := tx.Model(&RouteInfo{}).Where("route_token = ?", routeToken).First(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}

func UpdateRouteByToken(c *gin.Context, tx *gorm.DB, routeToken string, attrs map[string]interface{}) error {
	if tx == nil {
		tx = DB
	}
	r := tx.Model(&RouteInfo{}).Where("route_token = ?", routeToken).Update(attrs)
	return r.Error
}

func InsertRouteInfo(c *gin.Context, tx *gorm.DB, routeInfo *RouteInfo) (*RouteInfo, error) {
	err := Insert(tx, routeInfo)
	return routeInfo, err
}
