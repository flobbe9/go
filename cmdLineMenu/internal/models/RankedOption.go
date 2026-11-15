package models

import "iter"

type RankedOption struct {
	// the menu option getting a search ranking
	Option string

	// Number of points awarded depending on the searchQuery. Min 0.
	Points int
}

// Return a general purpose iterator.
func RankedOptionIterator(s []RankedOption) iter.Seq[RankedOption] {

	return func(yield func(r RankedOption) bool) {
		if len(s) == 0 {
			return;
		}

		for _, r := range s {
			if !yield(r) {
				return;
			}
		}
	}
}