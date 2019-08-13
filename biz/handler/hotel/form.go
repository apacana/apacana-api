package hotel

type GetAgodaHotelForm struct {
	HotelIDs []int64 `json:"hotel_ids" binding:"required"`
}
