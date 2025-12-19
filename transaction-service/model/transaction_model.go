package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	PaymentStatusPending = "pending"
	PaymentStatusSuccess = "success"
	PaymentStatusFailed  = "failed"
	PaymentStatusExpired = "expired"
	PaymentStatusCancel  = "cancel"
)

const (
	PaymentMethodQRIS = "qris"
)

const (
	FraudStatusAccept    = "accept"
	FraudStatusDeny      = "deny"
	FraudStatusChallange = "challange"
)

type Transaction struct {
	ID           uint   `json:"id" gorm:"primaryKey:autoIncrement"`
	Name         string `json:"name" gorm:"type:varchar(255);not null"`
	Phone        string `json:"phone" gorm:"type:varchar(20);not null"`
	Email        string `json:"email" gorm:"type:varchar(255)"`
	Address      string `json:"address" gorm:"type:text"`
	SubTotal     int64  `json:"sub_total" gorm:"type:bigint;not null"`
	TaxTotal 	 int64	`json:"tax_total" type:"bigint;not null"`
	GrandTotal   int64  `json:"grand_total" gorm:"type:bigint;not null"`
	MerchantID   uint   `json:"merchant_id" gorm:"type:bigint;not null"`
	

	//midtrans requirement
	PaymentStatus   string `json:"payment_status" gorm:"type:varchar(50);default:'pending'"`
	PaymentMethod   string `json:"payment_method" gorm:"type:varchar(50)"`
	PaymentCode     string `json:"payment_code" gorm:"type:varchar(100)"`
	OrderID         string `json:"order_id" gorm:"type:varchar(100);uniqueIndex"`
	TransactionCode string `json:"transaction_code" gorm:"type:varchar(100)"`
	PaymentToken    string `json:"payment_token" gorm:"type:text"`
	CallbackURL     string `json:"callback_url" gorm:"type:text"`
	ExpiredAt       *time.Time `json:"expired_at"`
	Notes 			string `json:"notes" gorm:"type:text"`
	Currency 		string `json:"currency" gorm:"type:varchar(10);default:'IDR'"`
	FraudStatus 	string `json:"fraud_status" gorm:"type:varchar(50)"`


	//virtual field for response
	MerchantName string `json:"merchant_name" gorm:"-"`


	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index  "`

	 TransactionProducts []TransactionProduct `gorm:"foreignKey:TransactionID"`
}