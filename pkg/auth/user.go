package auth

type User struct {
	ID     int64  `json:"id"`
	CityID int64  `json:"city_id"`
	Lang   string `json:"-"`
}
