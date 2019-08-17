package handler

import (
	"github.com/apacana/apacana-api/biz/handler/hotel"
	"github.com/apacana/apacana-api/biz/handler/stroke"
	"github.com/apacana/apacana-api/biz/handler/user"
	"github.com/apacana/apacana-api/biz/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	r.POST("/api/user/tourist/create/", user.CreateTourist)
	Api := r.Group("/api/")
	Api.Use(middleware.ApacanaCookieRequire)
	{
		ApiUser := Api.Group("/user/")
		{
			ApiUser.GET("/info/", user.GetUserInfo)
			ApiUser.POST("/register/", user.RegisterUser)
			ApiUser.POST("/login/", user.LoginUser)
		}
		ApiHotel := Api.Group("/hotel/")
		{
			ApiHotelAgoda := ApiHotel.Group("/agoda/")
			{
				ApiHotelAgoda.POST("/get/", hotel.GetAgodaHotel)
			}
		}
		ApiStroke := Api.Group("/stroke/")
		{
			ApiStroke.POST("create", stroke.CreateStroke)
		}
	}
}
