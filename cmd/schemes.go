// Copyright 2022-2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var (
	schemesCmd = &cobra.Command{
		Use:   "schemes",
		Short: "list the names of attestation schemes supported by the Veraison service",
		Long:  "List the names of attestation schemes supported by the Veraison service.",
		Args:  cobra.NoArgs,
		RunE:  doSchemesCommand,
	}
)

func doSchemesCommand(cmd *cobra.Command, args []string) error {
	schemes, err := service.GetSupportedSchemes()
	if err != nil {
		return err
	}

	sort.Strings(schemes)

	fmt.Printf("Attestation schemes supported by %s:\n", service.EndPointURI.Host)
	for _, scheme := range schemes {
		fmt.Println(scheme)
	}

	return nil
}
