package hotel

type GetAgodaHotelForm struct {
	HotelIDs []int64 `json:"hotel_ids" binding:"required"`
}

type SearchAgodaHotelForm struct {
	Latitude     string  `json:"latitude" binding:"required"`
	Longitude    string  `json:"longitude" binding:"required"`
	MaxResult    *string `json:"max_result" binding:"omitempty"`
	SearchRadius *string `json:"search_radius" binding:"omitempty"`
}
