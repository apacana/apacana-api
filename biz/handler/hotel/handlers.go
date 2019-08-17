package hotel

import (
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetAgodaHotel(c *gin.Context) {
	var getAgodaHotelForm GetAgodaHotelForm
	if err := c.ShouldBindJSON(&getAgodaHotelForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "GetAgodaHotel bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "GetAgodaHotel from: %+v", getAgodaHotelForm)

	hotelInfos, err := mysql.MGetHotelInfoAgodaByHotelID(c, nil, getAgodaHotelForm.HotelIDs)
	if err != nil && err != gorm.ErrRecordNotFound {
		helper.FormatLogPrint(helper.ERROR, "GetAgodaHotel MGetHotelInfoAgodaByHotelID failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
		return
	}
	data := map[int64]interface{}{}
	for _, hotel := range hotelInfos {
		data[hotel.HotelID] = hotel
	}
	helper.FormatLogPrint(helper.LOG, "GetAgodaHotel success, data: %+v", data)

	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, data)
}
