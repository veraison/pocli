// Copyright 2022-2023 Contributors to the Veraison project.
// SPDX-License-Identifier: Apache-2.0
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/veraison/apiclient/management"
)

var (
	getCmd = &cobra.Command{
		Use:   "get SCHEME [UUID]",
		Short: "get the active policy for the SCHEME, or, if the UUID is supplied, a specific policy associated with the scheme",
		Long: "Get the active policy for the SCHEME, or, if the UUID is supplied, " +
			"a specific policy associated with the scheme.\n" +
			"SCHEME must be a scheme name (use \"pocli schemes\" " +
			"to list valid values).\n" +
			"UUID is the unique identifier of the policy to be activated.",
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(1),
			cobra.MaximumNArgs(2),
			validateGetArgs,
		),
		RunE: doGetCommand,
	}

	getVersion        int32
	getOutputFilePath string
	getRulesFilePath  string
)

func init() {
	getCmd.PersistentFlags().StringVarP(&getOutputFilePath, "output", "o", "",
		"write the policy to the specified file, rather than STDOUT")
	getCmd.PersistentFlags().StringVarP(&getRulesFilePath, "write-rules", "w", "",
		"write the policy's rules to the specified file.")
}

func validateGetArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 2 {
		if _, err := uuid.Parse(args[1]); err != nil {
			return fmt.Errorf("invalid policy ID: %w", err)
		}
	}

	return nil
}

func doGetCommand(cmd *cobra.Command, args []string) error {
	var policy *management.Policy
	var err error
	var desc string

	if len(args) == 2 {
		desc = fmt.Sprintf("Policy %q:", args[1])
		policyID, err := uuid.Parse(args[1])
		if err != nil {
			return fmt.Errorf("invalid policy ID: %w", err)
		}

		policy, err = service.GetPolicy(args[0], policyID)
	} else {
		desc = "Active policy:"
		policy, err = service.GetActivePolicy(args[0])
	}

	if err != nil {
		return err
	}

	text, err := json.MarshalIndent(policy, "", "    ")
	if err != nil {
		return err
	}

	if getOutputFilePath == "" {
		fmt.Println(desc)
		fmt.Println(string(text))
	} else {
		if err := os.WriteFile(getOutputFilePath, text, 0o644); err != nil {
			return err
		}
	}

	if getRulesFilePath != "" {
		data := []byte(policy.Rules)
		if err := os.WriteFile(getRulesFilePath, data, 0o644); err != nil {
			return err
		}
	}

	return nil
}
