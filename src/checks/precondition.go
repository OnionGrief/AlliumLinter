package checks

import (
	"fmt"
	"refalLint/src/config"
	"refalLint/src/lexer"
	"refalLint/src/logger"
	"unicode"
)

// CheckPrecondition Проверки которые возможно выполнить до построения дерева
func CheckPrecondition(tokens []lexer.Token) []logger.Log {
	var logs []logger.Log
	logs = append(logs, checkExtern(tokens)...)
	if config.SnakeCase || config.CamelCase {
		logs = append(logs, checkNames(tokens)...)
	}
	logs = append(logs, checkConstCount(tokens)...)
	return logs
}

// Проверка на вынесение констант
func checkConstCount(tokens []lexer.Token) []logger.Log {
	var logs []logger.Log
	count := config.ConstCount
	if count < 3 {
		count = 3
	}
	leng := config.ConstLen
	if leng < 3 {
		leng = 3
	}
	chks := make(map[string]uint)
	str := ""
	for _, tkn := range tokens {
		switch tkn.TokenType {
		case lexer.ASCIIV:
			str += tkn.Value
		default:
			if str != "" {
				if len(str) > int(leng) {
					chks[str]++
					if chks[str] > count {
						logs = append(logs, logger.ExternLog(fmt.Sprintf("%s must be a const %d:%d", str, tkn.Start, tkn.Finish)))
					}
				}
				str = ""
			}
		}
	}
	return logs
}

// Проверяем на CamelCase или SnakeCase
func checkNames(tokens []lexer.Token) []logger.Log {
	var logs []logger.Log
	for _, tkn := range tokens {
		if tkn.TokenType != lexer.NAME && tkn.TokenType != lexer.VAR {
			continue
		}
		value := tkn.Value
		if tkn.TokenType == lexer.VAR {
			value = value[2:]
		}
		if config.CamelCase {
			if !isCamelCase(value) {
				logs = append(logs, logger.FormatNameLogCamel(tkn))
			}
		}
		if config.SnakeCase {
			if !isSnakeCase(value) {
				logs = append(logs, logger.FormatNameLogShake(tkn))
			}
		}
	}
	return logs
}
func checkExtern(tokens []lexer.Token) []logger.Log {
	var logs []logger.Log
	var extrn *lexer.Token
	for _, tkn := range tokens {
		if tkn.TokenType == lexer.EXTERN {
			if extrn == nil {
				extrn = &tkn
			} else {
				logs = append(logs, logger.ExternLog(fmt.Sprintf("Обнаружено несколько EXTERN, рекомендуется использовать 1 %d:%d (Первое использование %d:%d)", extrn.Start.Column, extrn.Start.Column, tkn.Start.Column, tkn.Start.Column)))
				return logs
			}
		}
	}
	return logs
}

func isSnakeCase(s string) bool {
	if len(s) == 0 {
		return false
	}

	for i, r := range s {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
		if r == '_' && (i == 0 || i == len(s)-1) {
			return false
		}
		if r == '_' && s[i-1] == '_' {
			return false
		}
	}
	return true
}

func isCamelCase(s string) bool {
	if len(s) == 0 {
		return false
	}

	foundUpper := false
	for _, r := range s {
		if unicode.IsUpper(r) {
			foundUpper = true
			continue
		}
		if !unicode.IsLower(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return foundUpper
}
