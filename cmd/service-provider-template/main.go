//go:generate opencontrolplane-gen
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openmcp-project/controller-utils/pkg/fips"

	// opencontrolplane-gen:replace github.com/openmcp-project/service-provider-template=MODULE template=PROVIDER_NAME
	"github.com/openmcp-project/service-provider-template/cmd/service-provider-template/app"
)

func main() {
	fips.Verify(context.Background())
	cmd := app.NewServiceServiceCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}
