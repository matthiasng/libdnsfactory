package main

import (
	"context"
	"fmt"

	"github.com/libdns/libdns"
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

	records, err := provider.AppendRecords(context.Background(), "test.com", []libdns.Record{
		{
			Type:  "TXT",
			Name:  "test.com",
			Value: "test",
			TTL:   120,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(records)
}
