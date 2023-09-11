package Model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type RuleForm struct {
	Product       string
	RuleContent   string
	ScheduledTime time.Time
}

type FlattenedRule struct {
	RuleForm
	FlattenedRule string
}
type AnalyzedRule struct {
	RuleForm
	InsightPool
	FlattenedRule string

	InsightTypes []InsightType
}

func (a *AnalyzedRule) Value() (driver.Value, error) {
	return json.Marshal(a.InsightTypes)
}

func (a *AnalyzedRule) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &a.InsightTypes)
}

type InsightType struct {
	InsightName   string
	InsightExists bool
}
type InsightPool struct {
	InsightTypes []InsightType
}
