package SIEM

import "gorm.io/gorm"

type rule struct {
	gorm.Model
	product string
	content string
}
