package Strategies

import (
	"awesomeProject1/SIEM/Model"
	"awesomeProject1/SIEM/splunk"
)

type SplunkStrategy struct{}

func (s *SplunkStrategy) InsightAnalysis(rule *Model.RuleForm) *Model.AnalyzedRule {

	analyzedRule := splunk.SetandExecuteInsightChain(rule)

	return analyzedRule
}
