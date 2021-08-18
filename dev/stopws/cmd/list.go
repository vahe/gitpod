// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package cmd

import (
	"github.com/gitpod-io/gitpod/stopws/pkg/stop"
	"github.com/spf13/cobra"
)

// listCmd lists affected pods
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists affected pods",
	Run: func(cmd *cobra.Command, args []string) {
		stop.ListPods()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
