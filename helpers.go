package main

import (
	"bytes"
	"log"
	"strings"
)

// Debugf prints out a debug log message; uses yellow color to indicate something
// other than normal output.
func Debugf(format string, a ...interface{}) {
	coloredFmt := strings.Join(
		[]string{
			"\033[33;1m<DEBUG> ",
			format,
			"\033[0m",
		}, "")
	log.Printf(coloredFmt, a...)
}

// Errorf prints out an error log message; uses red color to signal failure.
func Errorf(format string, a ...interface{}) {
	coloredFmt := strings.Join(
		[]string{
			"\033[41;1m<ERROR>\033[0m ",
			"\033[31;1m",
			format,
			"\033[0m",
		}, "")
	log.Printf(coloredFmt, a...)
}

// applyRule is just a helper function which is meant to make validation of
// rule logic easier by making it simple to pass any method through it as long
// as this method satisfies RuleTestFunc function prototype.
func applyRule(f RuleTestFunc, buf bytes.Buffer) (bool, error) {
	return f(buf)
}
