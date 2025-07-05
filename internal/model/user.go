package model

type User struct {
	ID        int
	Name      string
	HighScore int
	Coin      int
	Token     string
}

type UserCreateRequest struct {
	Name string `json:"name"`
}

type UserCreateResponse struct {
	Token string `json:"token"`
}
