package rule_insight

type Priority int

type Priorities []Priority

func (p Priorities) Len() int {
	return len(p)
}

func (p Priorities) Less(i, j int) bool {
	return p[i] < p[j]
}

func (p Priorities) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
