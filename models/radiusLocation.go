package models

type RadiusLocation struct {
	MaxLongitude     float64 `json:"maxLongitude"`
	MaxLatitude  float64 `json:"maxLatitude"`
	MinLongitude float64 `json:"minLongitude"`
	MinLatitude float64 `json:"minLatitude"`
}