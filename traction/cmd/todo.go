package cmd

import (
	"github.com/spf13/cobra"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage your todos from the command line",
}

func init() {
	todoCmd.AddCommand(todoListCmd)
}
