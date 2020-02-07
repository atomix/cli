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
	"context"
	"github.com/atomix/go-client/pkg/client/value"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

func newValueCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set {add,contains,remove,size,clear}",
		Short: "Manage the state of a distributed set",
	}
	addClientFlags(cmd)
	cmd.PersistentFlags().StringP("name", "n", "", "the set name")
	cmd.PersistentFlags().Lookup("name").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__atomix_get_sets"},
	}
	cmd.MarkPersistentFlagRequired("name")
	cmd.AddCommand(newValueSetCommand())
	cmd.AddCommand(newValueGetCommand())
	cmd.AddCommand(newValueWatchCommand())
	return cmd
}

func getValueName(cmd *cobra.Command) string {
	name, _ := cmd.Flags().GetString("name")
	return name
}

func getValue(cmd *cobra.Command, name string) value.Value {
	database := getDatabase(cmd)
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	m, err := database.GetValue(ctx, name)
	if err != nil {
		ExitWithError(ExitError, err)
	}
	return m
}

func newValueSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "set <value>",
		Args: cobra.ExactArgs(1),
		Run:  runValueSetCommand,
	}
	cmd.Flags().Uint64("version", 0, "the value version to update (for optimistic locking)")
	return cmd
}

func runValueSetCommand(cmd *cobra.Command, args []string) {
	val := getValue(cmd, getValueName(cmd))
	version, _ := cmd.Flags().GetUint64("version")
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	var newVersion uint64
	var err error
	if version > 0 {
		newVersion, err = val.Set(ctx, []byte(args[0]), value.IfVersion(version))
	} else {
		newVersion, err = val.Set(ctx, []byte(args[0]))
	}
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput(newVersion)
	}
}

func newValueGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "get",
		Args: cobra.NoArgs,
		Run:  runValueGetCommand,
	}
}

func runValueGetCommand(cmd *cobra.Command, _ []string) {
	val := getValue(cmd, getValueName(cmd))
	version, _ := cmd.Flags().GetUint64("version")
	ctx, cancel := getTimeoutContext(cmd)
	defer cancel()
	value, version, err := val.Get(ctx)
	if err != nil {
		ExitWithError(ExitError, err)
	} else {
		ExitWithOutput("Value: %v, Version: %d", value, version)
	}
}

func newValueWatchCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "watch",
		Args: cobra.NoArgs,
		Run:  runValueWatchCommand,
	}
}

func runValueWatchCommand(cmd *cobra.Command, _ []string) {
	val := getValue(cmd, getValueName(cmd))
	watchCh := make(chan *value.Event)
	if err := val.Watch(context.Background(), watchCh); err != nil {
		ExitWithError(ExitError, err)
	}

	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case event := <-watchCh:
			Output("Value: %v, Version: %d", event.Value, event.Version)
		case <-signalCh:
			ExitWithSuccess()
		}
	}
}