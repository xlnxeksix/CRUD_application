package splunk

type MacroEntry struct {
	Name    string            `json:"name"`
	Id      string            `json:"id"`
	Content MacroContentField `json:"content"`
}

type MacroContentField struct {
	Args       string `json:"args"`
	Definition string `json:"definition"`
}
