package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hedlx/doless/cli/ops"
	"github.com/spf13/cobra"
)

var runtimeCmd = &cobra.Command{
	Use:   "runtime",
	Short: "Runtime API methods",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var runtimeName string

var runtimeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runtime, err := ops.CreateRuntime(cmd.Context(), runtimeName, args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(runtime, "", "  ")
		fmt.Println(string(j))
	},
}

var runtimeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List",
	Run: func(cmd *cobra.Command, args []string) {
		runtimes, err := ops.ListRuntimes(cmd.Context())
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		j, _ := json.MarshalIndent(runtimes, "", "  ")
		fmt.Println(string(j))
	},
}

func init() {
	RootCmd.AddCommand(runtimeCmd)
	runtimeCmd.AddCommand(runtimeCreateCmd)
	runtimeCmd.AddCommand(runtimeListCmd)

	runtimeCreateCmd.Flags().StringVarP(&runtimeName, "name", "n", "", "name")
	runtimeCreateCmd.MarkFlagRequired("name")
}
