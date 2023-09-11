package rule_insight

import (
	"awesomeProject1/SIEM/Model"
)

type Insight interface {
	Execute(rule *Model.FlattenedRule, pool *Model.InsightPool)
	SetNext(Insight)
}