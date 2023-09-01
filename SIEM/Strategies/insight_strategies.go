package Strategies

import (
	"awesomeProject1/SIEM/Model"
)

type Insight interface {
	InsightAnalysis(ruleContent Model.RuleForm) Model.AnalyzedRule
}
