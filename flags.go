package main

import (
	"./cli"
	"flag"
)

func initializeFlags() *cli.Args {
	var a = &cli.Args{}

	flag.BoolVar(a.PTRDebug(), "debug", false, "Enable debug level log messages")
	flag.BoolVar(a.PTRpNamedSubexpr(), "print-named-subexp", false, "Prints out named sub-expressions, aka. named groups, JSON serialized")

	flag.BoolVar(a.PTRnoCheckNamedSubexpr(), "no-check-named-subexp", false, "If rule contains named sub-expressions don't validate them in matches")

	flag.BoolVar(a.PTRpExpandedTemplate(), "print-expanded-template", true, "If rule contains named sub-expressions, expand into the template if given")

	flag.StringVar(a.PTRrulesFilePath(), "rules-file", "rules", "Location of rules file")

	flag.Parse()
	return a
}
