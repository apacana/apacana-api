package mysql

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type HotelInfoAgoda struct {
	ID                  int64  `gorm:"id" json:"id"`
	HotelID             int64  `gorm:"hotel_id" json:"hotel_id"`
	ChainID             int64  `gorm:"chain_id" json:"chain_id"`
	ChainName           string `gorm:"chain_name" json:"chain_name"`
	BrandID             int64  `gorm:"brand_id" json:"brand_id"`
	BrandName           string `gorm:"brand_name" json:"brand_name"`
	HotelName           string `gorm:"hotel_name" json:"hotel_name"`
	HotelFormerlyName   string `gorm:"hotel_formerly_name" json:"hotel_formerly_name"`
	HotelTranslatedName string `gorm:"hotel_translated_name" json:"hotel_translated_name"`
	AddressLine1        string `gorm:"address_line1" json:"address_line1"`
	AddressLine2        string `gorm:"address_line2" json:"address_line2"`
	ZipCode             string `gorm:"zip_code" json:"zip_code"`
	City                string `gorm:"city" json:"city"`
	State               string `gorm:"state" json:"state"`
	Country             string `gorm:"country" json:"country"`
	CountryISOCode      string `gorm:"country_iso_code" json:"country_iso_code"`
	StarRating          string `gorm:"star_rating" json:"star_rating"`
	Longitude           string `gorm:"longitude" json:"longitude"`
	Latitude            string `gorm:"latitude" json:"latitude"`
	Url                 string `gorm:"url" json:"url"`
	CheckIn             string `gorm:"check_in" json:"check_in"`
	CheckOut            string `gorm:"check_out" json:"check_out"`
	NumberRooms         int64  `gorm:"number_rooms" json:"number_rooms"`
	NumberFloors        int64  `gorm:"number_floors" json:"number_floors"`
	YearOpened          int64  `gorm:"year_opened" json:"year_opened"`
	YearRenovated       int64  `gorm:"year_renovated" json:"year_renovated"`
	Photo1              string `gorm:"photo1" json:"photo1"`
	Photo2              string `gorm:"photo2" json:"photo2"`
	Photo3              string `gorm:"photo3" json:"photo3"`
	Photo4              string `gorm:"photo4" json:"photo4"`
	Photo5              string `gorm:"photo5" json:"photo5"`
	OverView            string `gorm:"over_view" json:"over_view"`
	RatesFrom           int64  `gorm:"rates_from" json:"rates_from"`
	ContinentID         int64  `gorm:"continent_id" json:"continent_id"`
	ContinentName       string `gorm:"continent_name" json:"continent_name"`
	CityID              int64  `gorm:"city_id" json:"city_id"`
	CountryID           int64  `gorm:"country_id" json:"country_id"`
	NumberOfReviews     int64  `gorm:"number_of_reviews" json:"number_of_reviews"`
	RatingAverage       string `gorm:"rating_average" json:"rating_average"`
	RatesCurrency       string `gorm:"rates_currency" json:"rates_currency"`
}

const (
	HotelInfoAgodaTableName = "hotel_info_agoda"
)

func (a *HotelInfoAgoda) TableName() string {
	return HotelInfoAgodaTableName
}

func MGetHotelInfoAgodaByHotelID(c *gin.Context, tx *gorm.DB, hotelIDs []int64) ([]*HotelInfoAgoda, error) {
	if tx == nil {
		tx = DB.Model(&HotelInfoAgoda{})
	}
	var ref []*HotelInfoAgoda
	r := tx.Where("hotel_id in (?)", hotelIDs).Find(&ref)
	if r.Error != nil {
		return nil, r.Error
	}
	return ref, nil
}
