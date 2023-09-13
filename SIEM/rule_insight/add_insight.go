package rule_insight

import (
	"awesomeProject1/SIEM/Model"
	"gorm.io/gorm"
)

func AddInsightsToDatabase(db *gorm.DB) {
	insight := Model.InsightType{
		InsightID:   1,
		InsightName: "Wildcard",
	}
	db.Create(&insight)
}
