package main

import (
	"./cli"
	"flag"
)

func initializeFlags() *cli.Args {
	var a = &cli.Args{}

	flag.BoolVar(a.PTRFailfast(), "failfast", false, "Fail after first unsuccessful rule")

	flag.BoolVar(a.PTRDebug(), "debug", false, "Enable debug level log messages")

	flag.BoolVar(a.PTRpNamedSubexpr(), "print-named-subexp", false, "Prints out named sub-expressions, aka. named groups, JSON serialized")

	flag.BoolVar(a.PTRnoCheckNamedSubexpr(), "no-check-named-subexp", false, "If rule contains named sub-expressions don't validate them in matches")

	flag.BoolVar(a.PTRnopExpandedTemplate(), "no-print-exp-template", false, "If rule contains named sub-expressions, do not expand into the template if given")

	flag.StringVar(a.PTRrulesFilePath(), "rules-file", "rules", "Location of rules file")

	flag.Parse()
	return a
}
