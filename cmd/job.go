package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Run and list jobs",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

var jobLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List jobs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from jobLsCmd func")
	},
}

var jobRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Define and run a new job",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hello from jobRunCmd func")
	},
}

func init() {
	jobCmd.AddCommand(jobLsCmd)
	jobCmd.AddCommand(jobRunCmd)
}
