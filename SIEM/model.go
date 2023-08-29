package SIEM

import "gorm.io/gorm"

type rule struct {
	gorm.Model
	SIEM        string
	RuleContent string
}
