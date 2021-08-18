// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/gitpod-io/gitpod/stopws/pkg/stop"
)

var singleCmd = &cobra.Command{
	Use:   "single",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("pod argument missing")
			return
		}
		pod := args[0]
		fmt.Println("pod", pod)
		stop.Single(pod)
	},
}

func init() {
	rootCmd.AddCommand(singleCmd)
}
