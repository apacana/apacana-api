package handler

import (
	"github.com/apacana/apacana-api/biz/handler/hotel"
	"github.com/apacana/apacana-api/biz/handler/point"
	"github.com/apacana/apacana-api/biz/handler/route"
	"github.com/apacana/apacana-api/biz/handler/stroke"
	"github.com/apacana/apacana-api/biz/handler/user"
	"github.com/apacana/apacana-api/biz/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.GET("/api/user/prepare/", user.PrepareUser) // 用户预准备
	Api := r.Group("/api/")
	Api.Use(middleware.ApacanaCookieRequire)
	{
		ApiUser := Api.Group("/user/")
		{
			ApiUser.GET("/info/", user.GetUserInfo)       // 获取用户信息
			ApiUser.POST("/register/", user.RegisterUser) // 用户注册
			ApiUser.POST("/login/", user.LoginUser)       // 用户登录
		}
		ApiHotel := Api.Group("/hotel/")
		{
			ApiHotelAgoda := ApiHotel.Group("/agoda/")
			{
				ApiHotelAgoda.POST("/get/", hotel.GetAgodaHotel) // 获得agoda酒店数据
			}
		}
		ApiStroke := Api.Group("/stroke/")
		{
			ApiStroke.POST("/create/", stroke.CreateStroke)          // 新建行程
			ApiStroke.GET("/:strokeToken/", stroke.GetStroke)        // 获取行程详细信息
			ApiStroke.POST("/change/default/", stroke.ChangeDefault) // 更改默认行程
		}
		ApiRoute := Api.Group("/route/")
		{
			ApiRoute.POST("/create/", route.CreateRoute)      // 新建路线
			ApiRoute.GET("/:routeToken/", route.GetRoute)     // 获取路线信息
			ApiRoute.POST("/close/", route.CloseRoute)        // 关闭路线关注
			ApiRoute.POST("/add_point/", route.AddRoutePoint) // 新增路线点
		}
		ApiPoint := Api.Group("/point/")
		{
			ApiPoint.POST("/add/", point.AddPoint)       // 新增行程点
			ApiPoint.POST("/delete/", point.DeletePoint) // 删除行程点
		}
	}
}
