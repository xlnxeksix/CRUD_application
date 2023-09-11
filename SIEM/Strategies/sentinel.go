package Strategies

import "awesomeProject1/SIEM/Model"

type SentinelStrategy struct{}

func (s *SentinelStrategy) InsightAnalysis(rule *Model.RuleForm) *Model.AnalyzedRule {
	var Arule *Model.AnalyzedRule
	return Arule
}
