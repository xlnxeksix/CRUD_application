package splunk

import (
	"awesomeProject1/SIEM/Model"
	"awesomeProject1/SIEM/rule_insight"
	"fmt"
	"gorm.io/gorm"
)

type WildCard struct {
	next rule_insight.Insight
	db   *gorm.DB
}

func (w *WildCard) Execute(r *Model.FlattenedRule, InsightIDs []int) {

	for i := 0; i < len(r.FlattenedRule); i++ {
		if r.FlattenedRule[i] == '*' {
			if i == 0 {
				InsightIDs = append(InsightIDs, 1)
				break
			}
			if r.RuleContent[i-1] != '\'' {
				InsightIDs = append(InsightIDs, 1)
				break
			}
		}
	}
	w.next.Execute(r, InsightIDs)
}

func (w *WildCard) SetNext(next rule_insight.Insight) {
	w.next = next
}

type Report struct {
	next rule_insight.Insight
}

func (rep *Report) Execute(r *Model.FlattenedRule, InsightIDs []int) {
	fmt.Println("Report exists")
}

func (rep *Report) SetNext(next rule_insight.Insight) {
	rep.next = next
}

func SetandExecuteInsightChain(rule *Model.RuleForm) []int {
	// Flatten the splunk macros
	var InsightIDs []int
	fmt.Println(FlattenQuery(rule.RuleContent))
	flattenedRule := &Model.FlattenedRule{
		RuleForm:      *rule,
		FlattenedRule: FlattenQuery(rule.RuleContent), // Assuming RuleContent is the flattened rule
	}
	wCard := &WildCard{}
	report := &Report{}

	wCard.SetNext(report)
	wCard.Execute(flattenedRule, InsightIDs)
	return InsightIDs
}
