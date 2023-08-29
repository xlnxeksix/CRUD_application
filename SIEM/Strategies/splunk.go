package Strategies

type SplunkStrategy struct{}

func (s *SplunkStrategy) InsightAnalysis(ruleContent string) string {

	return "Splunk" + ruleContent
}
