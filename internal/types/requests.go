package types

type UserLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateGift struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type UpdateGift struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type CreateOccasion struct {
	Name   string  `json:"name"`
	Guests []int64 `json:"guests"`
}

type UpdateOccasion struct {
	Name   string  `json:"name"`
	Guests []int64 `json:"guests"`
}
