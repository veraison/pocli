// Copyright 2022-2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create SCHEME RULES_FILE",
		Short: "create a new policy for SCHEME using the rules inside RULES_FILE",
		Long: "Create a new policy for SCHEME using the rules inside RULES_FILE.\n" +
			"SCHEME must be a scheme name (use \"pocli schemes\" " +
			"to list valid values).\n" +
			"RULES_FILE is the path to the file containing the rules for the policy.",
		Args: cobra.MatchAll(cobra.ExactArgs(2), validateCreateArgs),
		RunE: doCreateCommand,
	}

	createName    string
	doNotActivate bool
)

func init() {
	createCmd.PersistentFlags().StringVarP(&createName, "name", "n", "",
		"the name for the new policy")
	createCmd.PersistentFlags().BoolVarP(&doNotActivate, "dont-activate", "d", false,
		"if specified, the new policy will not be activated afer being created")
}

func validateCreateArgs(cmd *cobra.Command, args []string) error {
	if _, err := os.Stat(args[1]); err != nil {
		return fmt.Errorf("could not stat rules file: %w", err)
	}

	return nil
}

func doCreateCommand(cmd *cobra.Command, args []string) error {
	rules, err := os.ReadFile(args[1])
	if err != nil {
		return err
	}

	policy, err := service.CreateOPAPolicy(args[0], rules, createName)
	if err != nil {
		return err
	}

	fmt.Println("Policy created:")
	text, err := json.MarshalIndent(policy, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(text))

	if doNotActivate {
		return nil
	}

	err = service.ActivatePolicy(args[0], policy.UUID)
	if err != nil {
		return err
	}

	fmt.Println("Policy activated.")
	return nil
}
