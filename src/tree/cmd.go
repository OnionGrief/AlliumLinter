package tree

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"refalLint/src/lexer"
)

var TreeOut = &cobra.Command{
	Use:   "tree",
	Short: "Build tree for .ref file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("need only 1 file")
			return
		}
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Println("не смогли открыть файл")
			return
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		lex := lexer.NewLexer(reader)
		tkns, err := lex.GetAllTokens()
		if err != nil {
			fmt.Println(err)
			return
		}

		psr := NewParser(tkns)

		node, err := psr.Parse()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(CreateDot(node))

	},
}
