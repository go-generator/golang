package model

import "time"

const ()

type payment struct {
	Amount              int       `json:"amount" bson:"amount" gorm:"column:amount"`
	ChargeType          string    `json:"chargeType" bson:"chargeType" gorm:"column:chargeType"`
	CreditBankAccount   string    `json:"creditBankAccount" bson:"creditBankAccount" gorm:"column:creditBankAccount"`
	CurrencyCode        string    `json:"currencyCode" bson:"currencyCode" gorm:"column:currencyCode"`
	DebitBankAccount    string    `json:"debitBankAccount" bson:"debitBankAccount" gorm:"column:debitBankAccount"`
	ExternalSystemId    string    `json:"externalSystemId" bson:"externalSystemId" gorm:"column:externalSystemId"`
	FeeAmount           int       `json:"feeAmount" bson:"feeAmount" gorm:"column:feeAmount"`
	FeeDebitBankAccount string    `json:"feeDebitBankAccount" bson:"feeDebitBankAccount" gorm:"column:feeDebitBankAccount"`
	FeeDueDate          time.Time `json:"feeDueDate" bson:"feeDueDate" gorm:"column:feeDueDate"`
	PayeeId             string    `json:"payeeId" bson:"payeeId" gorm:"column:payeeId"`
	PayerId             string    `json:"payerId" bson:"payerId" gorm:"column:payerId"`
	PaymentDate         time.Time `json:"paymentDate" bson:"paymentDate" gorm:"column:paymentDate"`
	PaymentId           string    `json:"paymentId" bson:"_id" gorm:"column:_id:primary_key"`
	PaymentStatus       string    `json:"paymentStatus" bson:"paymentStatus" gorm:"column:paymentStatus"`
	PostingType         string    `json:"postingType" bson:"postingType" gorm:"column:postingType"`
	TransFeeId          int       `json:"transFeeId" bson:"transFeeId" gorm:"column:transFeeId"`
	UpdatedDate         time.Time `json:"updatedDate" bson:"updatedDate" gorm:"column:updatedDate"`
}
