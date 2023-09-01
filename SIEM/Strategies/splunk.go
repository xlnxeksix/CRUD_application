package Strategies

import (
	"awesomeProject1/SIEM/Model"
	Insight2 "awesomeProject1/SIEM/insights/splunk_insight"
	"fmt"
)

type SplunkStrategy struct{}

func (s *SplunkStrategy) InsightAnalysis(rule Model.RuleForm) Model.AnalyzedRule {
	var Arule Model.AnalyzedRule
	Arule.Product = rule.Product
	Arule.RuleContent = rule.RuleContent
	Insight2.SetandExecuteInsightChain(&Arule)
	fmt.Println(Arule.WildCardInsight)
	return Arule
}
