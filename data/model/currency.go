package Model

import (
	"github.com/jinzhu/gorm"
)

//modeling table currency
type Currency struct {
	gorm.Model
	From string `gorm:"type:varchar(100);"`
	To   string `gorm:"type:varchar(100);"`
	Rate []Rate
}
