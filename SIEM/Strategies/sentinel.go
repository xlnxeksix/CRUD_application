package Strategies

type SentinelStrategy struct{}

func (s *SentinelStrategy) InsightAnalysis(ruleContent string) string {

	return "Sentinel" + ruleContent
}
