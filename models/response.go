package models

type Response struct{
	Success bool `json:"success"`
	//error bool `json.error`
	Message string `json:"message"`
	Error error `json:"error"`
	Data []LocationData `json:"data"`
	DataLength int `json:"dataLength"`
}