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
	listCmd = &cobra.Command{
		Use:   "list SCHEME",
		Short: "list policies for the specified SCHEME",
		Long: "List policies for the specified SCHEME.\n" +
			"SCHEME must be a scheme name (use \"pocli schemes\" to list valid values).",
		Args: cobra.ExactArgs(1),
		RunE: doListCommand,
	}

	listName           string
	listOutputFilePath string
)

func init() {
	listCmd.PersistentFlags().StringVarP(&listName, "name", "n", "",
		"if specified, only policies with the specified name will be listed")
	listCmd.PersistentFlags().StringVarP(&listOutputFilePath, "output", "o", "",
		"write the policies to the specified file, rather than STDOUT")
}

func doListCommand(cmd *cobra.Command, args []string) error {
	policies, err := service.GetPolicies(args[0], listName)
	if err != nil {
		return err
	}

	text, err := json.MarshalIndent(policies, "", "    ")
	if err != nil {
		return err
	}

	if listOutputFilePath == "" {
		if len(policies) > 0 {
			if listName == "" {
				fmt.Printf("Policies for scheme %s:\n", args[0])
			} else {
				fmt.Printf("Policies for scheme %s with name %s:\n",
					args[0], listName)
			}

			fmt.Println(string(text))
		} else { // zero policies returned
			if listName == "" {
				fmt.Printf("No policies for scheme %s.\n", args[0])
			} else {
				fmt.Printf("No policies for scheme %s with name %s.\n",
					args[0], listName)
			}
		}
	} else {
		if err := os.WriteFile(listOutputFilePath, text, 0o644); err != nil {
			return err
		}
	}

	return nil
}
