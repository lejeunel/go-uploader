package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Run and list jobs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from jobCmd func")
	},
}

var jobLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List jobs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from jobLsCmd func")
	},
}

func init() {

	jobCmd.AddCommand(jobLsCmd)
}
