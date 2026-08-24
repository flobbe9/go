package models;

type Options struct {
	// Whether to erase all menu lines (that is all selectOptions) after user submit. 
	IsClearMenuOnSubmit bool;

	// Whether to print a hint next to the question telling the user how to submit.
	IsShowSubmitHint bool;
	
	// Whether to print the selected answer next to the question.
	IsDisplayAnswer bool;
}