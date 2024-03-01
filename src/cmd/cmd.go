package cmd

import (
	"github.com/OnionGrief/AlliumLinter/src/cmd/lint"
	"github.com/OnionGrief/AlliumLinter/src/tree"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "refLint",
	Short: "A Refal linter application",
}

func Execute() {

	RootCmd.AddCommand(lint.LintCmd)
	RootCmd.AddCommand(tree.TreeOut)
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	// Инициализация глобальных флагов и других настроек
}
