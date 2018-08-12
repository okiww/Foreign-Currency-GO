package Model

import (
	"time"

	"github.com/jinzhu/gorm"
)

//modeling table currency
type Rate struct {
	gorm.Model
	Date       time.Time
	CurrencyID uint
	Currency   Currency `gorm:"foreignkey:CurrencyID"`
	Rate       float64  `gorm:"size:3;"`
}
