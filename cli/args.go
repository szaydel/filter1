package cli

// Args is a globally accessible command line flags (arguments).
type Args struct {
	debug               bool
	noCheckNamedSubexpr bool
	nopExpandedTemplate bool
	pNamedSubexpr       bool
	rulesFilePath       string
}

func (a *Args) PTRDebug() *bool {
	return &a.debug
}

func (a *Args) PTRpNamedSubexpr() *bool {
	return &a.pNamedSubexpr
}

func (a *Args) PTRrulesFilePath() *string {
	return &a.rulesFilePath
}

func (a *Args) PTRnopExpandedTemplate() *bool {
	return &a.nopExpandedTemplate
}

func (a *Args) PTRnoCheckNamedSubexpr() *bool {
	return &a.noCheckNamedSubexpr
}

func (a Args) Debug() bool {
	return a.debug
}

func (a Args) RulesFilePath() string {
	return a.rulesFilePath
}

func (a Args) PrintNamedSubexpr() bool {
	return a.pNamedSubexpr
}

func (a Args) NoPrintExpandedTemplate() bool {
	return a.nopExpandedTemplate
}
