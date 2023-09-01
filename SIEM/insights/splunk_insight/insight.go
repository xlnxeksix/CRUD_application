package Insight

import (
	"awesomeProject1/SIEM/Model"
)

type Analysis interface {
	Execute(rule *Model.AnalyzedRule)
	SetNext(Analysis)
}

func SetandExecuteInsightChain(rule *Model.AnalyzedRule) *Model.AnalyzedRule {
	wCard := &WildCard{}
	report := &Report{}

	wCard.SetNext(report)
	wCard.Execute(rule)

	return rule
}
