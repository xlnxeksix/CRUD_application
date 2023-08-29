package Strategies

type SplunkCloudStrategy struct{}

func (s *SplunkCloudStrategy) InsightAnalysis(ruleContent string) string {

	return "SplunkCloud" + ruleContent
}
