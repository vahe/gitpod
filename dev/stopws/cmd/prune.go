// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package cmd

import (
	"github.com/gitpod-io/gitpod/stopws/pkg/prune"
	"github.com/spf13/cobra"
)

var pruneCmd = &cobra.Command{
	Use:   "prune",
	Run: func(cmd *cobra.Command, args []string) {
		prune.PrunePods()
	},
}

func init() {
	rootCmd.AddCommand(pruneCmd)
}
