package Model

type RuleForm struct {
	Product     string
	RuleContent string
}

type AnalyzedRule struct {
	RuleForm
	WildCardInsight bool
}
