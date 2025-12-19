package model

import "time"

type Warehouse struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name" gorm:"varchar(100);notnull"`
	Address   string     `json:"address" gorm:"type-text"`
	Photo     string     `json:"photos" gorm:"type-text"`
	Phone     string     `json:"phone" gorm:"varchar(100);notnull"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	WarehouseProducts []WarehouseProduct `json:"warehouse_products" gorm:"foreignKey:WarehouseID"`
}
