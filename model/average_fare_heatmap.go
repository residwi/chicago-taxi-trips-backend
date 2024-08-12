package model

type AverageFareHeatmap struct {
	S2ID        string  `json:"s2id"`
	AverageFare float64 `json:"fare"`
}

type FarePerLocation struct {
	Latitude  float64
	Longitude float64
	Fare      float64
}
