package models

type Response struct{
	Success bool `json:"success"`
	//error bool `json.error`
	Message string `json:"message"`
	Data []LocationData `json:"data"`
	DataLength int `json:"dataLength"`
}