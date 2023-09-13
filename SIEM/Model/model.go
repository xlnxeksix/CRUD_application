package Model

import (
	"gorm.io/gorm"
	"time"
)

type RuleForm struct {
	Product       string
	RuleContent   string
	ScheduledTime time.Time
}

type FlattenedRule struct {
	gorm.Model
	RuleForm
	FlattenedRule string
	Insights      []*InsightType `gorm:"many2many:insight_table;"`
}

type InsightType struct {
	gorm.Model
	InsightID   int
	InsightName string
}
