# Linode Metadata Client for Go

This package contains an incomplete client to interact 
with the Linode Metadata service in Go.

## Usage

NOTE: Since this package is not currently public, you will need to use the
[replace directive](https://go.dev/ref/mod#go-mod-file-replace) to reference a
local version of this repository.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/linode/linodego/metadata"
)

func main() {
	// Create and initialize the metadata client.
	// By default, this function automatically generates a token.
	// This functionality can be disabled using the opts argument.
	client, err := metadata.NewClient(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Get info from /v1/instance
	instanceInfo, err := client.GetInstance(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Get info from /v1/network
	networkInfo, err := client.GetNetwork(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Get info from /v1/ssh-keys
	sshKeys, err := client.GetSSHKeys(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Get info from /v1/user-data
	userData, err := client.GetUserData(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Print various info
	fmt.Println("Instance ID:", instanceInfo.InstanceID)
	fmt.Println("Instance Public IP:", networkInfo.IPv4.Public[0])
	fmt.Println("Authorized SSH Keys:", sshKeys.Users.Root)
	fmt.Println("Authorized Users:", sshKeys.Users.Username)
	fmt.Println("User Data:", userData)
}
```