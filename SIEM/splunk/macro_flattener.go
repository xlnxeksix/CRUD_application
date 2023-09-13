package splunk

import (
	"fmt"
	"regexp"
	"strings"
)

/*
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
*/
func FlattenQuery(q string) string {
	// Define your macros within the function
	macros := map[string]MacroEntry{
		"simple_macro": {
			Name: "simple_macro",
			Content: MacroContentField{
				Definition: "index=8",
			},
		},
		"circular_macro_a": {
			Name: "circular_macro_a",
			Content: MacroContentField{
				Definition: "index=9 and `circular_macro_b`",
			},
		},
		"circular_macro_b": {
			Name: "circular_macro_b",
			Content: MacroContentField{
				Definition: "index=10 and `circular_macro_a`",
			},
		},
		"macro_with_params(1)": {
			Name: "macro_with_params",
			Content: MacroContentField{
				Definition: "index=$value1$",
				Args:       "value1",
			},
		},
		"search_head_cluster": {
			Name: "search_head_cluster",
			Content: MacroContentField{
				Definition: "searching",
				Args:       "value1",
			},
		},
	}

	// Create a closure function for recursive expansion
	expandedMacros := make(map[string]bool)

	// Create a closure function for recursive expansion
	var expandMacro func(string) string
	expandMacro = func(input string) string {
		reg := regexp.MustCompile("`(([\\w|-]+)(\\(.+?\\))?)`")
		reMatches := reg.FindAllStringSubmatch(input, -1)

		// Replace macros with their expanded content
		for _, match := range reMatches {
			macroName := match[1]
			if macroEntry, ok := macros[macroName]; ok {
				// Check if the macro has already been expanded to prevent circular references
				if expandedMacros[macroName] {
					// Handle circular reference, e.g., by returning an error or breaking the loop
					return "Circular reference detected!"
				}

				// Mark the macro as expanded
				expandedMacros[macroName] = true

				macroContent := macroEntry.Content.Definition

				// Parse and process macro arguments
				if macroEntry.Content.Args != "" {
					argMacro := fmt.Sprintf("$%s$", macroEntry.Content.Args)
					argValue := match[3] // The argument value is captured by the third submatch

					// Remove parentheses from the argument value
					argValue = strings.TrimPrefix(argValue, "(")
					argValue = strings.TrimSuffix(argValue, ")")

					macroContent = strings.ReplaceAll(macroContent, argMacro, argValue)
				}
				// Recursively expand nested macros
				macroContent = expandMacro(macroContent)

				input = strings.ReplaceAll(input, match[0], macroContent)
			}
		}

		return input
	}

	// Use the closure function to expand macros
	return expandMacro(q)
}
