package SIEM

import (
	"awesomeProject1/SIEM/Model"
	"gorm.io/gorm"
)

type SIEMRepository interface {
	InsertRule(rule *Model.AnalyzedRule) error
}

type SQLRuleRepository struct {
	DB *gorm.DB
}

func NewSQLRuleRepository(db *gorm.DB) SIEMRepository {
	return &SQLRuleRepository{DB: db}
}

func (repo *SQLRuleRepository) InsertRule(rule *Model.AnalyzedRule) error {
	insertQuery := "INSERT INTO analyzed_rules (product, rule_content, ScheduledTime, InsightTypes, FlattenedRule) VALUES (?, ?, ?)"
	if err := repo.DB.Exec(insertQuery, rule.Product, rule.RuleContent, rule.ScheduledTime, rule.InsightTypes, rule.FlattenedRule).Error; err != nil {
		return err
	}
	return nil
}
