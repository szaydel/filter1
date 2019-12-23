package main

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
)

// Rules is simply a list of Rule(s) structures defined above.
type Rules struct {
	list []Rule
}

// RuleByIndex returns a rule at a given index, or a nil if index is invalid.
func (rules *Rules) RuleByIndex(idx int) *Rule {
	if idx < len(rules.list) {
		return &rules.list[idx]
	}
	return nil
}

// Len returns number of rules in the rules list.
func (rules *Rules) Len() int {
	return len(rules.list)
}

// ListOfRules returns a slice of Rule(s) structs.
func (rules *Rules) ListOfRules() []Rule {
	return rules.list
}

// LoadRules consumes rules from an io.Reader, assumes data to be JSON encoded,
// and attempts to json.Unmarshal(...) this data into a slice of rules.
func (rules *Rules) LoadRules(r io.Reader) (int, error) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r)
	if err != nil {
		return -1, err
		// log.Printf("Error (LoadRules): %v", err)
	}
	if err := json.Unmarshal(buf.Bytes(), &rules.list); err != nil {
		return -1, err
	}

	return rules.Len(), nil
}

// CompileRules receives an io.Reader and from this reader loads rules which
// are subsequently used by the regexp matching code.
func CompileRules(r io.Reader) *CompiledRuleList {
	var rules = &Rules{}
	rules.LoadRules(r)
	var crl CompiledRuleList

	// Build a slice of all compiled patterns which we then use with data
	// consumed from stdin.
	var cexp *regexp.Regexp
	var pattern string
	for i, r := range rules.ListOfRules() {
		if r.PatString == MatchNoOutput {
			pattern = zeroInputPattern
		} else {
			pattern = r.PatString
		}
		if r.ERESyntax {
			cexp = regexp.MustCompilePOSIX(pattern)
		} else {
			cexp = regexp.MustCompile(pattern)
		}
		crl.Append(CompiledRule{re: cexp, desc: rules.RuleByIndex(i)})
	}

	return &crl
}
