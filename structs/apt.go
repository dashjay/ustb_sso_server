package structs

import "fmt"

type AppToken struct {
	ID        string `json:"id"`
	AppURL    string `json:"app_url"`
	APPID     string `json:"appid"`
	ReturnURL string `json:"return_url"`
	RandToken string `json:"rand_token"`
}

func (v *AppToken) Encode() string {
	return fmt.Sprintf("appid=%s&return_url=%s&rand_token=%s", v.APPID, v.ReturnURL, v.RandToken)
}

type State struct {
	State int    `json:"state"`
	Data  string `json:"data"`
}
