module github.com/thomas-maurice/ringos/ringo-cli

go 1.17

require (
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/thomas-maurice/ringos/go-client v0.0.0-20220306114144-ab298e74e094
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
)

replace github.com/thomas-maurice/ringos/go-client => ../go-client
