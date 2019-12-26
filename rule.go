package main

// Rule provides a description, or a specification for exactly what the
// given rule is going to match, and what the various constraints are.
type Rule struct {
	Comment       string            `json:"comment"`
	PatString     string            `json:"pat_string"`
	Expect        string            `json:"expect"`
	Occurences    int               `json:"occurences"`
	Constraints   Constraints       `json:"constraints"`
	ERESyntax     bool              `json:"ere_syntax"`
	NoCheckNSE    bool              `json:"no_check_named_subexprs"`
	Debug         bool              `json:"debug"`
	NamedSubExprs map[string]string `json:"named_subexprs"`
	Template      string            `json:"template"`
}
