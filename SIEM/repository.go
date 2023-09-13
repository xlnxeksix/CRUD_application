package SIEM

import (
	"awesomeProject1/SIEM/Model"
	"gorm.io/gorm"
)

type SIEMRepository interface {
	InsertRule(rule *Model.RuleForm) error
	InsertInsight(rule *Model.FlattenedRule, InsightID []int) error
}

type SQLRuleRepository struct {
	DB *gorm.DB
}

func NewSQLRuleRepository(db *gorm.DB) SIEMRepository {
	return &SQLRuleRepository{DB: db}
}

func (repo *SQLRuleRepository) InsertRule(rule *Model.RuleForm) error {
	insertQuery := "INSERT INTO rule_form (product, rule_content, ScheduledTime) VALUES (?, ?, ?)"
	if err := repo.DB.Exec(insertQuery, rule.Product, rule.RuleContent, rule.ScheduledTime).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SQLRuleRepository) InsertInsight(rule *Model.FlattenedRule, InsightID []int) error {
	insights := make([]*Model.InsightType, 0) // Create a slice of InsightType pointers

	// Retrieve each InsightType using its ID and append to the insights slice
	for _, id := range InsightID {
		insight := &Model.InsightType{}
		if err := repo.DB.First(insight, id).Error; err != nil {
			return err
		}
		insights = append(insights, insight)
	}

	// Append all retrieved insights to the InsightTypes association of the rule
	if err := repo.DB.Model(rule).Association("InsightTypes").Append(insights); err != nil {
		return err
	}

	return nil
}
