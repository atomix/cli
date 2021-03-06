// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"fmt"
	"github.com/atomix/go-client/pkg/client/counter"
	"github.com/spf13/cobra"
	"strconv"
)

func newCounterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "counter <name> [...]",
		Short:              "Manage the state of a distributed counter",
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If only the name was specified, open an interactive shell
			name := args[0]
			if name == "-h" || name == "--help" {
				return cmd.Help()
			}
			if len(args) == 1 {
				ctx := getContext()
				if ctx == nil {
					ctx = newContext("atomix", "counter", name)
					setContext(ctx)
				} else {
					ctx = ctx.withCommand("counter", name)
				}
				return ctx.run()
			}

			// Get the command for the specified operation
			var subCmd *cobra.Command
			op := args[1]
			switch op {
			case "get":
				subCmd = newCounterGetCommand(name)
			case "set":
				subCmd = newCounterSetCommand(name)
			case "increment":
				subCmd = newCounterIncrementCommand(name)
			case "decrement":
				subCmd = newCounterDecrementCommand(name)
			case "help", "-h", "--help":
				if len(args) == 2 {
					helpCmd := &cobra.Command{
						Use:   fmt.Sprintf("counter %s [...]", name),
						Short: "Manage the state of a distributed counter",
					}
					helpCmd.AddCommand(newCounterGetCommand(name))
					helpCmd.AddCommand(newCounterSetCommand(name))
					helpCmd.AddCommand(newCounterIncrementCommand(name))
					helpCmd.AddCommand(newCounterDecrementCommand(name))
					return helpCmd.Help()
				} else {
					var helpCmd *cobra.Command
					switch args[2] {
					case "get":
						helpCmd = newCounterGetCommand(name)
					case "set":
						helpCmd = newCounterSetCommand(name)
					case "increment":
						helpCmd = newCounterIncrementCommand(name)
					case "decrement":
						helpCmd = newCounterDecrementCommand(name)
					default:
						return fmt.Errorf("unknown command %s", args[2])
					}
					return helpCmd.Help()
				}
			default:
				return fmt.Errorf("unknown command %s", op)
			}
			addClientFlags(subCmd)

			// Set the arguments after the name and execute the command
			subCmd.SetArgs(args[2:])
			return subCmd.Execute()
		},
	}
	return cmd
}

func getCounter(cmd *cobra.Command, name string) (counter.Counter, error) {
	database, err := getDatabase(cmd)
	if err != nil {
		return nil, err
	}
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	return database.GetCounter(ctx, name)
}

func newCounterGetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			counter, err := getCounter(cmd, name)
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := counter.Get(ctx)
			if err != nil {
				return err
			}
			cmd.Println(value)
			return nil
		},
	}
	return cmd
}

func newCounterSetCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "set <value>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			counter, err := getCounter(cmd, name)
			if err != nil {
				return err
			}
			value, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			err = counter.Set(ctx, int64(value))
			if err != nil {
				return err
			}
			cmd.Println(value)
			return nil
		},
	}
	return cmd
}

func newCounterIncrementCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "increment [delta]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			counter, err := getCounter(cmd, name)
			if err != nil {
				return err
			}
			var delta int64 = 1
			if len(args) > 0 {
				value, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				delta = int64(value)
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := counter.Increment(ctx, delta)
			if err != nil {
				return err
			}
			cmd.Println(value)
			return nil
		},
	}
	return cmd
}

func newCounterDecrementCommand(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "decrement [delta]",
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			counter, err := getCounter(cmd, name)
			if err != nil {
				return err
			}
			var delta int64 = 1
			if len(args) > 0 {
				value, err := strconv.Atoi(args[0])
				if err != nil {
					return err
				}
				delta = int64(value)
			}
			ctx, cancel := getTimeoutContext(cmd)
			defer cancel()
			value, err := counter.Decrement(ctx, delta)
			if err != nil {
				return err
			}
			cmd.Println(value)
			return nil
		},
	}
	return cmd
}
