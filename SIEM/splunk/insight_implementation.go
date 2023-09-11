package splunk

import (
	"awesomeProject1/SIEM/Model"
	"awesomeProject1/SIEM/rule_insight"
	"fmt"
)

type WildCard struct {
	next rule_insight.Insight
}

func (w *WildCard) Execute(r *Model.FlattenedRule, pool *Model.InsightPool) {

	for i := 0; i < len(r.FlattenedRule); i++ {
		if r.FlattenedRule[i] == '*' {
			if i == 0 {
				pool.InsightTypes = append(pool.InsightTypes, Model.InsightType{"Wildcard", true})
				break
			}
			if r.RuleContent[i-1] != '\'' {
				pool.InsightTypes = append(pool.InsightTypes, Model.InsightType{"Wildcard", true})
				break
			}
		}
	}
	pool.InsightTypes = append(pool.InsightTypes, Model.InsightType{"Wildcard", false})
	w.next.Execute(r, pool)
}

func (w *WildCard) SetNext(next rule_insight.Insight) {
	w.next = next
}

type Report struct {
	next rule_insight.Insight
}

func (rep *Report) Execute(r *Model.FlattenedRule, pool *Model.InsightPool) {
	fmt.Println("Report exists")
}

func (rep *Report) SetNext(next rule_insight.Insight) {
	rep.next = next
}

func SetandExecuteInsightChain(rule *Model.RuleForm) *Model.AnalyzedRule {
	// Flatten the splunk macros
	flattenedRule := &Model.FlattenedRule{
		RuleForm:      *rule,
		FlattenedRule: rule.RuleContent, // Assuming RuleContent is the flattened rule
	}

	// Create a new InsightPool
	pool := &Model.InsightPool{} // Initialize the pool

	wCard := &WildCard{}
	report := &Report{}

	wCard.SetNext(report)
	wCard.Execute(flattenedRule, pool)

	// Create an AnalyzedRule object
	analyzedRule := &Model.AnalyzedRule{
		RuleForm:      *rule,
		InsightPool:   *pool,                       // Dereference pool to include its content
		FlattenedRule: flattenedRule.FlattenedRule, // Use the flattened rule from flattenedRule
	}

	return analyzedRule
}
