package hotel

import (
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
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

func SearchHotel(c *gin.Context) {
	var searchAgodaHotelForm SearchAgodaHotelForm
	if err := c.ShouldBindJSON(&searchAgodaHotelForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "SearchHotel bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "SearchHotel from: %+v", searchAgodaHotelForm)

	lat := searchAgodaHotelForm.Latitude
	lng := searchAgodaHotelForm.Longitude
	maxResult := "30"
	searchRadius := "1"
	if searchAgodaHotelForm.MaxResult != nil {
		maxResult = *searchAgodaHotelForm.MaxResult
	}
	if searchAgodaHotelForm.SearchRadius != nil {
		searchRadius = *searchAgodaHotelForm.SearchRadius
	}

	requestUrl := "https://sherpa.agoda.com/Affiliate/FetchHotelsV2?" +
		"refKey=WHmVY%2BEwqwInzKYu3N907g%3D%3D&" +
		"apiRequest.criteria.additional.language=zh-cn&" +
		"apiRequest.criteria.additional.discountOnly=false&" +
		"apiRequest.criteria.additional.currency=CNY&" +
		"apiRequest.criteria.checkInDate=2019-12-09&" +
		"apiRequest.criteria.checkOutDate=2019-12-10&" +
		"apiRequest.criteria.geo.latitude=" + lat + "&" +
		"apiRequest.criteria.geo.longitude=" + lng + "&" +
		"apiRequest.criteria.geo.searchRadius=" + searchRadius + "&" +
		"apiRequest.criteria.additional.occupancy.numberOfRoom=1&" +
		"apiRequest.criteria.additional.occupancy.numberOfAdult=2&" +
		"apiRequest.criteria.additional.occupancy.numberOfChildren=0&" +
		"apiRequest.criteria.additional.maxResult=" + maxResult + "&" +
		"cid=1831827"
	resp, err := http.Get(requestUrl)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "SearchHotel close http body failed, err: %v", err)
		}
	}()
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "SearchHotel http request failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "SearchHotel http read body failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, string(body))
}
