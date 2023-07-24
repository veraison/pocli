// Copyright 2022-2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	activateCmd = &cobra.Command{
		Use:   "activate SCHEME UUID",
		Short: "activate policy UUID for scheme SCHEME",
		Long: "Activate policy UUID for scheme SCHEME.\n" +
			"SCHEME must be a scheme name (use \"pocli schemes\" " +
			"to list valid values).\n" +
			"UUID is the unique identifier of the policy to be activated.",
		Args: cobra.MatchAll(cobra.ExactArgs(2), validateActivateArgs),
		RunE: doActivateCommand,
	}
)

func validateActivateArgs(cmd *cobra.Command, args []string) error {
	if _, err := uuid.Parse(args[1]); err != nil {
		return fmt.Errorf("invalid policy ID: %w", err)
	}

	return nil
}

func doActivateCommand(cmd *cobra.Command, args []string) error {
	policyID, err := uuid.Parse(args[1])
	if err != nil {
		return fmt.Errorf("invalid policy ID: %w", err)
	}

	if err := service.ActivatePolicy(args[0], policyID); err != nil {
		return err
	}

	fmt.Println("Policy activated.")
	return nil
}
