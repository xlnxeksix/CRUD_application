package splunk

import (
	"awesomeProject1/SIEM/Model"
	"fmt"
	"regexp"
)

func FindMacroMatches(r *Model.RuleForm) {

	reg := regexp.MustCompile("`(([\\w|-]+)(\\(.+?\\))?)`")
	reMatches := reg.FindAllStringSubmatch(r.RuleContent, -1)
	var macros []string
	for _, match := range reMatches {
		macros = append(macros, match[1])
	}

	// Print the extracted macros
	for _, macro := range macros {
		fmt.Println(macro)
	}
}
