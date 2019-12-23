package main

const (
	// Exactly : exactly this many matches
	Exactly Constraints = iota
	// AtLeast : at least this many matches, or more
	AtLeast
	// AtMost : at most this many matches, or fewer
	AtMost

	unlimitedMatches = -1
	zeroInputPattern = `^$`
	MatchNoOutput    = "$EMPTY$"

	failedToApplyFmt   = "\033[1;36m<%s>\033[0m \033[1;33m%s\033[0m failed to apply"
	errApplyRuleMsgFmt = "🚩 \033[48;5;52mcheck[%s]:\033[0m \033[1;31m%v\033[0m"
)
