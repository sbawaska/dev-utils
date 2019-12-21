module github.com/projectriff/developer-utils

go 1.13

require (
	github.com/Azure/go-autorest v13.3.1+incompatible // indirect
	github.com/projectriff/stream-client-go v0.0.0-20191115170130-f3286a439f7b
	github.com/spf13/cobra v0.0.5
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v0.0.0-20190819141724-e14f31a72a77
)

replace github.com/projectriff/stream-client-go v0.0.0-20191115170130-f3286a439f7b => github.com/sbawaska/stream-client-go v0.0.0-20191220235430-889c18474076
