package main

import (
	"bytes"
	"fmt"
	// "log"
	"regexp"
)

// CompiledRule is a result of converting individual rules in the supplied
// rules file into compiled regular expressions.
type CompiledRule struct {
	re    *regexp.Regexp
	desc  *Rule
	debug bool
}

// Apply uses regexp engine to match given data received via buffer against
// compiled regexp with specifications in the rule. The rule is assumed
// to apply iff all tests are successful. For example, if a rule contains named
// sub-expressions, also frequently called named capture groups, the values for
// these named capture groups specified in the rule must match those in
func (c *CompiledRule) Apply(buf bytes.Buffer) (bool, map[string]error) {
	return c.applyRule(buf)
}

// Debug returns a boolean indicating whether or not debugging is enabled for
// this rule.
func (c *CompiledRule) Debug() bool {
	return c.debug
}

// Expand performs template rendering by taking values extracted with named
// sub-expressions and replacing variables with names matching those of
// named sub-expressions. For example, this template: "Hi, my name is $first
// $last, I am pleased to meet you.", will be successfully expanded if named
// sub-expressions contain groups named first and last. In this example, if
// first = "alpha" and last = "omega", the resulting byte slice will contain
// string: "Hi, my name is alpha omega, I am pleased to meet you."
func (c *CompiledRule) Expand(buf bytes.Buffer) []byte {
	if c.re.NumSubexp() > 0 {
		return c.expandTemplate(buf)
	}
	return []byte{}
}

// NamedSubExprsAsMap returns a map of named sub-expressions, also frequently
// called named capture groups. Return type is a map[string][]byte, where the
// keys are names of the capture group, and the values are actual extracted
// sub-string matches.
func (c *CompiledRule) NamedSubExprsAsMap(buf bytes.Buffer) map[string][]byte {
	if c.re.NumSubexp() > 0 {
		return c.nse(buf)
	}
	return map[string][]byte{}
}

func (c *CompiledRule) applyRule(buf bytes.Buffer) (bool, map[string]error) {
	var errs = make(map[string]error)
	var failed bool
	var testFns = make(map[string]RuleTestFunc)

	if c.desc.FailIfMatch {
		testFns["Must not Match"] = c.testMustNotMatch
	} else {
		testFns["Number of Matches"] = c.testMatchCount
	}

	// If we did not want to skip checking Named sub-expressions, and there is
	// at least one named sub-expression defined in the rule, add test to map
	// of tests.
	if c.desc.NoCheckNSE == false && len(c.desc.NamedSubExprs) > 0 {
		testFns["Named Sub-expr"] = c.testNamedSubexprs
	}

	for name, fn := range testFns {
		if ok, err := applyRuleTestFunc(fn, buf); !ok {
			failed = true
			errs[name] = err
		}
	}

	return !failed, errs
}

func (c *CompiledRule) testMustNotMatch(buf bytes.Buffer) (bool, error) {
	var matches [][]byte = c.re.FindAll(buf.Bytes(), unlimitedMatches)
	var got = len(matches)
	if got > 0 {
		return false, fmt.Errorf("%s:  <%s>, %d matches found", c.desc.Comment, c.re.String(), got)
	}
	return true, nil
}

func (c *CompiledRule) testMatchCount(buf bytes.Buffer) (bool, error) {
	var matches [][]byte = c.re.FindAll(buf.Bytes(), unlimitedMatches)
	var got = len(matches)

	// log.Printf("matches: %v", c.NamedSubExprsAsMap(buf))

	switch c.desc.Constraints {
	case Exactly:
		if got != c.desc.Occurences {
			return false, fmt.Errorf("%s: Number of occurences is incorrect; expected: (exactly %d), got %d", c.desc.Comment, c.desc.Occurences, got)
		}
	case AtLeast:
		if got < c.desc.Occurences {
			return false, fmt.Errorf("%s: Number of occurences is incorrect; expected: (at least %d), got: %d", c.desc.Comment, c.desc.Occurences, got)
		}
	case AtMost:
		if got > c.desc.Occurences {
			return false, fmt.Errorf("%s: Number of occurences is incorrect; expected (at most %d), got: %d", c.desc.Comment, c.desc.Occurences, got)
		}
	}
	return true, nil
}

func (c *CompiledRule) nse(buf bytes.Buffer) map[string][]byte {
	var namedSubexprs = make(map[string][]byte)
	var subexpNames = c.re.SubexpNames()
	var submatches = c.re.FindAllSubmatch(buf.Bytes(), unlimitedMatches)

	for i, lines := range submatches {
		for j, elem := range lines {
			if c.Debug() {
				Debugf("submatches[%d][%d] => %s", i, j, elem)
			}
			if j == 0 {
				continue
			}
			if _, ok := namedSubexprs[subexpNames[j]]; ok {
				key := fmt.Sprintf("%s_%d", subexpNames[j], i)
				namedSubexprs[key] = elem
			} else {
				namedSubexprs[subexpNames[j]] = elem
			}
		}
	}
	if c.Debug() {
		Debugf("%#v", namedSubexprs)
	}
	return namedSubexprs
}

func (c *CompiledRule) testNamedSubexprs(buf bytes.Buffer) (bool, error) {
	// If POSIX ERE syntax was chosen, we do not even have support for named
	// capture groups. If this syntax was chosen, ignore named groups even if
	// for whatever reason rule has them configured.
	if c.desc.ERESyntax {
		return true, nil
	}
	var namedSubexprs = c.nse(buf)

	// If rule contains an object with named subexpressions, this object is
	// effectively a map of subexpression names to their expected values.
	// When this object is present in the rule, we assume that extracted
	// sub-expressions will have same value as described in the rule. For
	// example, if rule has an object {"alpha": "omega"}, we expect that the
	// regexp which we compile and use to match and extract sub-strings will
	// have named capture group called "alpha", and the value in the matched
	// string is going to be "omega".
	for name, expected := range c.desc.NamedSubExprs {
		if _, ok := namedSubexprs[name]; !ok {
			return false,
				fmt.Errorf("RegEx pattern has no sub-expression named <%s>", name)
		}
		got := string(namedSubexprs[name])
		if got != expected {
			return false, fmt.Errorf("%s: Sub-expression named <%s> does not have expected value; expected: %s, got: %s", c.desc.Comment, name, expected, got)
		}
	}
	return true, nil
}

func (c *CompiledRule) expandTemplate(buf bytes.Buffer) []byte {
	template := []byte(c.desc.Template)
	result := []byte{}

	// For each match of the regex in the content.
	for _, submatches := range c.re.FindAllSubmatchIndex(buf.Bytes(), unlimitedMatches) {
		// Apply the captured submatches to the template and append the output
		// to the result.
		result = c.re.Expand(result, template, buf.Bytes(), submatches)
	}
	return result
}
