package models

type RankedOption struct {
	// the menu option getting a search ranking
	Option string;

	// Number of points awarded depending on the searchQuery. Min 0.
	Points int;
}