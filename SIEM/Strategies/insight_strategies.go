package Strategies

type Insight interface {
	InsightAnalysis(ruleContent string) string
}
