package cli

// Args is a globally accessible command line flags (arguments).
type Args struct {
	debug               bool
	noCheckNamedSubexpr bool
	pExpandedTemplate   bool
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

func (a *Args) PTRpExpandedTemplate() *bool {
	return &a.pExpandedTemplate
}

func (a *Args) PTRnoCheckNamedSubexpr() *bool {
	return &a.pExpandedTemplate
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

func (a Args) PrintExpandedTemplate() bool {
	return a.pExpandedTemplate
}
