package types

type CreateGiftResponse struct {
	Id int64 `json:"id"`
}

type GiftResponse struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Id   int64  `json:"id"`
}

type UsersResponse struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
}

type OccasionResponse struct {
	Name   string          `json:"name"`
	Id     int64           `json:"id"`
	Guests []UsersResponse `json:"guests"`
}

type CreateOccasionResponse struct {
	Id int64 `json:"id"`
}
