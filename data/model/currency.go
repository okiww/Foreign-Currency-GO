package Model

import (
	"time"

	"github.com/jinzhu/gorm"
)

//modeling table currency
type Currency struct {
	gorm.Model
	Date *time.Time
	From string  `gorm:"type:varchar(100);"`
	To   string  `gorm:"type:varchar(100);"`
	Rate float32 `gorm:"size:100"` // set field size to 255
}
