package Strategies

type QradarStrategy struct{}

func (s *QradarStrategy) InsightAnalysis(ruleContent string) string {

	return "Qradar" + ruleContent
}
