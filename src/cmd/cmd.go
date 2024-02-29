package cmd

import (
	"github.com/spf13/cobra"
	"refalLint/src/cmd/lint"
	"refalLint/src/tree"
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
