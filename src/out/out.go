package out

import "fmt"

const yellowColorCode = "\033[33m"
const redColorCode = "\033[31m"
const resetColorCode = "\033[0m"

func Warning(text string) {
	fmt.Println(yellowColorCode + text + resetColorCode)
}
func Error(text string) {
	fmt.Println(redColorCode + text + resetColorCode)
}
