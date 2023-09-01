package Insight

import (
	"awesomeProject1/SIEM/Model"
	"fmt"
)

type Report struct {
	next Analysis
}

func (rep *Report) Execute(r *Model.AnalyzedRule) {
	fmt.Println("Report exists")
}

func (rep *Report) SetNext(next Analysis) {
	rep.next = next
}
