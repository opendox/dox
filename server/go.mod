module github.com/opendox/dox/server

go 1.25.6

require (
	github.com/opendox/dox/packages/shared v0.0.0
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)

replace github.com/opendox/dox/packages/shared => ../packages/shared
