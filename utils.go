package wings

import (
	"github.com/agnivade/levenshtein"
)

// FindLevenshtein finds a target in a list of strings using levenshtein distance.
func FindLevenshtein(target string, list []string) (string, int) {
	var s string
	md := 100000
	for _, name := range list {
		d := levenshtein.ComputeDistance(target, name)
		if d < md {
			md = d
			s = name
		}
	}
	return s, md
}
