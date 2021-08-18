// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package cmd

import (
	"fmt"

	"github.com/gitpod-io/gitpod/stopws/pkg/stop"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use: "all",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("all")
		stop.All()
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
}
