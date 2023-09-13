package Strategies

import (
	"awesomeProject1/SIEM/Model"
	"awesomeProject1/SIEM/splunk"
)

type SplunkStrategy struct{}

func (s *SplunkStrategy) InsightAnalysis(rule *Model.RuleForm) []int {
	return splunk.SetandExecuteInsightChain(rule)
}
