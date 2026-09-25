# Usage 
```
options := []string{"apple", "banana", "ananas"};
selection, err := cmdLineMenu.Prompt(
    "Select a fruit",
    options,
    models.Config{IsClearMenuOnSubmit: true, IsShowSubmitHint: false, IsDisplayAnswer: true, IsOptionSearchEnabled: true},
)
// handle err
``` 

## Multi-select
```
options := []string{"apple", "banana", "ananas"};
var selections []string;

fmt.Println("Select a fruit");
for {
    selection, err := cmdLineMenu.Prompt(
        "", // leave empty for proper menu display
        options,
        models.Config{IsClearMenuOnSubmit: true, IsShowSubmitHint: false, IsDisplayAnswer: true, IsOptionSearchEnabled: true},
    )
    // handle err

    selections = append(selections, selection);
    deleteIndex := slices.Index(options, selection);
    options = slices.Delete(options, deleteIndex, deleteIndex + 1);
}
``` 