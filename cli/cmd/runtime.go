package cmd

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hedlx/doless/cli/ops"
	"github.com/hedlx/doless/cli/tui/runtime"
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

type runtimeCL struct {
	ctx context.Context
}

func (r *runtimeCL) Create(name string, path string) tea.Cmd {
	return func() tea.Msg {
		rt, err := ops.CreateRuntime(r.ctx, name, path)

		return runtime.RuntimeCreateResponseMsg{
			Resp: &runtime.RuntimeCreateResponse{
				Runtime: rt,
				Err:     err,
			},
		}
	}
}

func (r *runtimeCL) List() tea.Cmd {
	return func() tea.Msg {
		rt, err := ops.ListRuntimes(r.ctx)

		return runtime.RuntimeListResponseMsg{
			Resp: &runtime.RuntimeListResponse{
				Runtimes: rt,
				Err:      err,
			},
		}
	}
}

var runtimeCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		m := &runtime.RuntimeCreateModel{
			Name: runtimeName,
			Path: args[0],
			Creator: &runtimeCL{
				ctx: cmd.Context(),
			},
		}
		p := tea.NewProgram(runtime.InitRuntimeCreateModel(m))

		if err := p.Start(); err != nil {
			fmt.Printf("Error: %s", err)
		}
	},
}

var runtimeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List",
	Run: func(cmd *cobra.Command, args []string) {
		m := &runtime.RuntimeListModel{
			Lister: &runtimeCL{
				ctx: cmd.Context(),
			},
		}
		p := tea.NewProgram(runtime.InitRuntimeListModel(m))

		if err := p.Start(); err != nil {
			fmt.Printf("Error: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(runtimeCmd)
	runtimeCmd.AddCommand(runtimeCreateCmd)
	runtimeCmd.AddCommand(runtimeListCmd)

	runtimeCreateCmd.Flags().StringVarP(&runtimeName, "name", "n", "", "name")
}
