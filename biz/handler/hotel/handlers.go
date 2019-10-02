package hotel

import (
	"fmt"
	"github.com/apacana/apacana-api/biz/dal/mysql"
	"github.com/apacana/apacana-api/biz/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
		"apiRequest.criteria.checkInDate=" + searchAgodaHotelForm.CheckInData + "&" +
		"apiRequest.criteria.checkOutDate=" + searchAgodaHotelForm.CheckOutData + "&" +
		"apiRequest.criteria.geo.latitude=" + searchAgodaHotelForm.Latitude + "&" +
		"apiRequest.criteria.geo.longitude=" + searchAgodaHotelForm.Longitude + "&" +
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

func SearchHotelBooking(c *gin.Context) {
	var agodaHotelBookingForm AgodaHotelBookingForm
	if err := c.ShouldBindJSON(&agodaHotelBookingForm); err != nil {
		helper.FormatLogPrint(helper.WARNING, "SearchHotel bind json failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeParmErr, nil)
		return
	}
	helper.FormatLogPrint(helper.LOG, "SearchHotelBooking from: %+v", agodaHotelBookingForm)

	body := "{" +
		"\"criteria\": {" +
		"\"additional\": {" +
		"\"currency\": \"CNY\"," +
		"\"discountOnly\": false," +
		"\"language\": \"zh-cn\"," +
		"\"occupancy\": {" +
		"\"numberOfAdult\": 2," +
		"\"numberOfChildren\": 0}}," +
		"\"checkInDate\": \"" + agodaHotelBookingForm.CheckInDate + "\"," +
		"\"checkOutDate\": \"" + agodaHotelBookingForm.CheckOutDate + "\"," +
		"\"hotelId\": ["
	for _, hotelID := range agodaHotelBookingForm.HotelIDs {
		body += strconv.FormatInt(hotelID, 10)
	}
	body += "]}}"
	bodyReader := strings.NewReader(body)
	req, err := http.NewRequest("POST", "http://affiliateapi7643.agoda.com/affiliateservice/lt_v1", bodyReader)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "1831827:836ddc58-5e3a-4ffd-ad52-f5154315c7df")
	clt := http.Client{}
	resp, err := clt.Do(req)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			helper.FormatLogPrint(helper.ERROR, "SearchHotelBooking close http body failed, err: %v", err)
		}
	}()
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "SearchHotelBooking http request failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helper.FormatLogPrint(helper.ERROR, "SearchHotelBooking http read body failed, err: %v", err)
		helper.BizResponse(c, http.StatusOK, helper.CodeFailed, nil)
	}
	helper.BizResponse(c, http.StatusOK, helper.CodeSuccess, string(respBody))
}
