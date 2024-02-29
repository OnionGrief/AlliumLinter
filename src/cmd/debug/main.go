package main

import (
	"bufio"
	"fmt"
	"os"
	"refalLint/src/lexer"
	"refalLint/src/tree"
)

// Пакет используется строго для дебага дерева, не собирайте его

func main() {
	// Проверка, был ли предоставлен путь к файлу в аргументах командной строки
	if len(os.Args) < 2 {
		fmt.Println("Пожалуйста, укажите путь к файлу")
		return
	}
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
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
	//fmt.Println(tkns)
	//fmt.Println(checks.CheckPrecondition(tkns))
	psr := tree.NewParser(tkns)

	node, _ := psr.Parse()
	node = node
	psr.Checks()
	//fmt.Println(tree.CreateDot(node))
}
