package model

type Batch struct {
	BatchId   string `bson:"_id" json:"batchId"`
	BatchName string `bson:"batchName" json:"batchName"`
}
