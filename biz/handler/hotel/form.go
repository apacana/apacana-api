package hotel

type GetAgodaHotelForm struct {
	HotelIDs []int64 `json:"hotel_ids" binding:"required"`
}

type SearchAgodaHotelForm struct {
	Latitude     string  `json:"latitude" binding:"required"`
	Longitude    string  `json:"longitude" binding:"required"`
	CheckInData  string  `json:"check_in_data" binding:"required"`
	CheckOutData string  `json:"check_out_data" binding:"required"`
	Adult        *string `json:"adult" binding:"omitempty"`
	MaxResult    *string `json:"max_result" binding:"omitempty"`
	SearchRadius *string `json:"search_radius" binding:"omitempty"`
}

type AgodaHotelBookingForm struct {
	HotelIDs     []int64 `json:"hotel_ids" binding:"required"`
	CheckInDate  string  `json:"check_in_date" binding:"required"`
	CheckOutDate string  `json:"check_out_date" binding:"required"`
}
