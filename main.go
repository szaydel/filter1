package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	var rulesFile *os.File
	var openErr error
	var config = initializeFlags() // process CLI arguments

	if rulesFile, openErr = os.Open(config.RulesFilePath()); openErr != nil {
		panic(openErr)
	}

	var listOfCompiled = CompileRules(rulesFile)
	var buf bytes.Buffer
	var expandedTemplates bytes.Buffer
	var n int64
	var err error
	
	n, err = buf.ReadFrom(os.Stdin)
	if err != nil {
		panic(err)
	}
	if config.Debug() {
		Debugf("Read %d bytes from stdin", n)
	}
	var retcode int
	var pairs = make(map[string]string)

	for _, r := range listOfCompiled.Slice() {
		ok, errs := r.Apply(buf)
		if !ok {
			retcode = 1
			log.Printf(failedToApplyFmt, r.desc.Comment, r.re.String())
			for name, e := range errs {
				log.Printf(errApplyRuleMsgFmt, name, e)
			}
			// if ok...
		} else {
			// if we want to dump named sub-expressions...
			if config.PrintNamedSubexpr() {
				for k, v := range r.NamedSubExprsAsMap(buf) {
					pairs[k] = string(v)
				}
			}
			// when a template is present in the rule and there are named
			// sub-expressions in the pattern, replace associated variables
			// in template with sub-expressions captured by this regex pattern.
			if r.desc.Template != "" && r.re.NumSubexp() > 0 {
				expandedTemplates.Write(r.expandTemplate(buf))
			}
		}
	}
	// Once we are done processing all the rules, if the pairs map is not empty,
	// we dump its contents as a JSON object, to facilitate consumption with
	// tools such as `jq`.
	if config.PrintNamedSubexpr() && len(pairs) > 0 {
		if encoded, err := json.Marshal(pairs); err != nil {
			Errorf("%v", err)
		} else {
			fmt.Printf("%s\n", encoded)
		}
	}

	if config.PrintExpandedTemplate() {
		fmt.Printf("%s\n", expandedTemplates.String())
	}

	os.Exit(retcode)
}
