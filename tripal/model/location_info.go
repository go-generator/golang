package model

type LocationInfo struct {
	LocationInfoID string `json:"_id" gorm:"primary_key;column:_id" bson:"_id"`
	Rate           int32  `json:"rate" bson:"rate"`
	Rate1          int32  `json:"rate1" bson:"rate1" `
	Rate2          int32  `json:"rate2" bson:"rate2"`
	Rate3          int32  `json:"rate3" bson:"rate3"`
	Rate4          int32  `json:"rate4" bson:"rate4"`
	Rate5          int32  `json:"rate5" bson:"rate5"`
	RateLocation   int32  `json:"rateLocation" bson:"rateLocation"`
}

func (LocationInfo) CollectionName() string {
	return "locationInfo"
}
