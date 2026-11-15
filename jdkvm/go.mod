module jdkvm

go 1.26.2

require (
	github.com/spf13/cobra v1.10.2
	golang.org/x/sys v0.50.0
)

require github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203 // indirect

require (
	github.com/flobbe9/go/utils v0.3.0
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)

// TODO remove
replace github.com/flobbe9/go/utils => ../utils
