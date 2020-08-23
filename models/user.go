package models

type User struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}
