package getopt

import (
	"os"
	"slices"
	"strings"

	"github.com/lucasepe/x/getopt"
)

func EnvOrOptVal(envKey string, opts []getopt.OptArg, lookup []string) (val string) {
	val = os.Getenv(envKey)
	if val != "" {
		return
	}

	return OptVal(opts, lookup)
}

func OptVal(opts []getopt.OptArg, lookup []string) (val string) {
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

// AllOptArgs restituisce tutti gli argomenti associati a una lista di opzioni.
func AllOptArgs(opts []getopt.OptArg, keys []string) []string {
	result := []string{}

	keySet := make(map[string]struct{}, len(keys))
	for _, key := range keys {
		keySet[key] = struct{}{}
	}

	for _, opt := range opts {
		if _, ok := keySet[opt.Option]; ok {
			result = append(result, opt.Argument)
		}
	}
	return result
}
