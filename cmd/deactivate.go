// Copyright 2022-2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	deactivateCmd = &cobra.Command{
		Use:   "deactivate SCHEME",
		Short: "deactivate all policies for scheme SCHEME",
		Long: "Deactivate all policies for scheme SCHEME.\n" +
			"SCHEME must be a scheme name (use \"pocli schemes\" " +
			"to list valid values).\n",
		Args: cobra.ExactArgs(1),
		RunE: doDeactivateCommand,
	}
)

func doDeactivateCommand(cmd *cobra.Command, args []string) error {
	if err := service.DeactivateAllPolicies(args[0]); err != nil {
		return err
	}

	fmt.Printf("All policies for scheme %s deactivated.\n", args[0])
	return nil
}
