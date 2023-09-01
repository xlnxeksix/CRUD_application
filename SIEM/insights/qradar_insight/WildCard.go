package Insight

import (
	"awesomeProject1/SIEM/Model"
	"fmt"
)

type WildCard struct {
	next Analysis
}

func (w *WildCard) Execute(r *Model.AnalyzedRule) {

	fmt.Println("Wildcard analysis is carried out")
	for i := 0; i < len(r.RuleContent); i++ {
		if r.RuleContent[i] == '*' {
			if i == 0 {
				fmt.Println("Wildcard exists")
				r.WildCardInsight = true
				break
			}
			if r.RuleContent[i-1] != '\'' {
				fmt.Println("Wildcard exists")
				r.WildCardInsight = true
				break
			}
		}
	}
	w.next.Execute(r)
}

func (w *WildCard) SetNext(next Analysis) {
	w.next = next
}
