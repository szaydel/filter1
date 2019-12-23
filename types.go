package main

import "bytes"

// Constraints is an enum describing constraints for a given rule. The intent
// here is to be able to tell the matching logic whether we want to precisely
// match XX occurences, i.e. actual == rule, or to match at least XX occurences,// i.e. actual >= rule or to match at most XX occurences, i.e. actual <= rule.
// at most XX occurences.
type Constraints int

// RuleTestFunc is a type used for a function whose purpose is to return a
// boolean and an error after validating a given rule. Implementation details
// of this function are irrelevant, as long as on success it returns (true, nil)
// and on failure it returns (false, error).
type RuleTestFunc func(buf bytes.Buffer) (bool, error)
