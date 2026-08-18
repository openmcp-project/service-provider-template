//go:generate opencontrolplane-gen
package main

import (
	"fmt"
	"os"

	// opencontrolplane-gen:replace github.com/openmcp-project/service-provider-template=MODULE
	"github.com/openmcp-project/service-provider-template/cmd/service-provider-template/app"
)

func main() {
	cmd := app.NewServiceServiceCommand()

	if err := cmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
