package lint

import (
	"fmt"
	"github.com/OnionGrief/AlliumLinter/src/config"
	"github.com/OnionGrief/AlliumLinter/src/linter"
	"github.com/OnionGrief/AlliumLinter/src/out"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var fromDirectory string
var configPath string

var LintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint .ref files",
	Run: func(cmd *cobra.Command, args []string) {
		var files []string
		if fromDirectory != "" {
			files = findRefFilesInDirectory(fromDirectory)
		} else {
			if !allFilesRef(args) {
				log.Fatal("all files must be .ref")
			}
			files = args
		}
		cfg := config.ReadConfigFromFile(configPath)
		if cfg != nil {
			config.Cfg = *cfg
		}
		prepareFiles(files)
	},
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
func allFilesRef(slice []string) bool {
	for _, str := range slice {
		if len(str) < 4 || str[len(str)-4:] != ".ref" {
			return false
		}
		if !fileExists(str) {
			log.Fatalf("file %s not found", str)
		}
	}
	return true
}
func findRefFilesInDirectory(dir string) []string {
	var files []string
	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return files
	}

	for _, file := range fileList {
		fullPath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			files = append(files, findRefFilesInDirectory(fullPath)...)
		} else if strings.HasSuffix(file.Name(), ".ref") {
			files = append(files, fullPath)
		}
	}

	return files
}
func init() {
	LintCmd.Flags().StringVarP(&fromDirectory, "from-dir", "d", "", "Specify the directory to search for .ref files")
	LintCmd.Flags().StringVarP(&configPath, "cfg-path", "p", "./config.yaml", "Config file from thos directory")
	LintCmd.Flags().BoolVarP(&config.Cfg.SnakeCase, "snake", "s", false, "Use SnakeCase for format")
	LintCmd.Flags().BoolVarP(&config.Cfg.CamelCase, "camel", "c", false, "Use CamelCase for format")
	LintCmd.Flags().UintVarP(&config.Cfg.ConstLen, "constLen", "L", 3, "ConstLen")
	LintCmd.Flags().UintVarP(&config.Cfg.ConstCount, "constCount", "C", 3, "ConstCount")
	LintCmd.Flags().UintVarP(&config.Cfg.BlockLen, "blockLen", "b", 3, "Длина переиспользуемого блока")
}

func prepareFiles(files []string) {
	fmt.Println("Полученные файлы", files)
	results := make(chan linter.Result, len(files))
	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, filePath := range files {
		go linter.LintFile(filePath, results, &wg)
	}
	wg.Wait()
	close(results)

	outLogsWithFile(results)
}
func outLogsWithFile(result chan linter.Result) {

	for val := range result {
		if val.Success {
			for _, msg := range val.Logs {
				out.Warning(fmt.Sprintf("[%s] %s\n", val.FilePath, msg))
			}
		} else {
			out.Error(fmt.Sprintf("[%s] %s\n", val.FilePath, val.Error.Error()))
		}
	}
}
