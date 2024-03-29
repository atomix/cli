// SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package docs

import "github.com/spf13/cobra"

func GetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate-docs",
		Aliases: []string{"gen-docs"},
		Hidden:  true,
		RunE:    run,
	}
	cmd.Flags().StringP("output", "o", ".", "the path to which to output the docs")
	cmd.Flags().Bool("markdown", false, "generate docs in markdown format")
	cmd.Flags().Bool("man", false, "generate docs in man format")
	cmd.Flags().Bool("yaml", false, "generate docs in YAML format")
	_ = cmd.MarkFlagDirname("output")
	return cmd
}
