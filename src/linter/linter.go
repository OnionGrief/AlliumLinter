package linter

import (
	"bufio"
	"fmt"
	"os"
	"refalLint/src/checks"
	"refalLint/src/lexer"
	"refalLint/src/logger"
	"refalLint/src/tree"
	"sync"
	"time"
)

type Result struct {
	FilePath string
	Success  bool
	Error    error
	Logs     []logger.Log
}

// LintFile асинхронно проверяет файл и возвращает результат через канал
func LintFile(filePath string, results chan<- Result, wg *sync.WaitGroup) {
	var out []logger.Log
	defer wg.Done() // Уменьшаем счетчик WaitGroup по завершению горутины

	file, err := os.Open(filePath)
	if err != nil {
		results <- Result{FilePath: filePath, Error: fmt.Errorf("не смогли открыть файл")}
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lex := lexer.NewLexer(reader)
	tkns, err := lex.GetAllTokens()
	if err != nil {
		results <- Result{FilePath: filePath, Error: err}
		return
	}
	fmt.Println("get all tokens ok")
	logs := checks.CheckPrecondition(tkns)
	out = append(out, logs...)
	psr := tree.NewParser(tkns)
	node, err := psr.Parse()
	if err != nil {
		results <- Result{FilePath: filePath, Error: err}
		return
	}
	logs2 := psr.Checks()
	out = append(out, logs2...)
	node = node
	// После проверки отправляем результат в канал
	results <- Result{FilePath: filePath, Success: true, Logs: out}
}
func runWithTimeout(fn func(), timeout time.Duration) error {
	done := make(chan error, 1)

	go func() {
		defer close(done)
		fn()
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("выполнение заняло более %f секунд, возможно произошла ошибка при разборе", timeout.Seconds())
	}
}
