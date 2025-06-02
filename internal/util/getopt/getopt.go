package getopt

import (
	"slices"
	"strings"

	"github.com/lucasepe/x/getopt"
)

func FindOptVal(opts []getopt.OptArg, lookup []string) (val string) {
	for _, opt := range opts {
		if slices.Contains(lookup, opt.Opt()) {
			val = opt.Argument
			break
		}
	}

	return
}

func HasOpt(opts []getopt.OptArg, lookup []string) bool {
	for _, opt := range opts {
		if slices.Contains(lookup, opt.Opt()) {
			return true
		}
	}

	return false
}

func WantsHelp(args []string) bool {
	if len(args) == 0 {
		return true
	}

	return strings.EqualFold(args[0], "help")
}
