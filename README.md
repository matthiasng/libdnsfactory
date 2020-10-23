libdnsfactory
=============
[![godoc reference](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/matthiasng/libdnsfactory)

With `libdnsfactory`, you can create libdns providers by name.

This allows you to support all providers in a dynamic way. For example, if you want to let the user decide which provider to use with which configuration.

The [factory](https://github.com/matthiasng/libdnsfactory/blob/master/factory.go) and [docs](https://github.com/matthiasng/libdnsfactory/blob/master/docs.md) are [generated](https://github.com/matthiasng/libdnsfactory/blob/master/generate/main.go) from the provider repositories.

# Install

```sh
go get -u github.com/matthiasng/libdnsfactory
```


# Example

```go
package main

import (
	"fmt"

	"github.com/matthiasng/libdnsfactory"
)

func main() {
    name := "hetzner"
    configMap := map[string]string{
		"AuthAPIToken": "<your token>",
    }
    
	provider, err := libdnsfactory.NewProvider(name, configMap)
	if err != nil {
		panic(err)
	}

    records, err := provider.AppendRecords(...)
    if err != nil {
		panic(err)
    }
    
    fmt.Println(records)
}

```

