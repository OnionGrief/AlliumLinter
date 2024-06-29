package checks

import (
	"fmt"
	"github.com/OnionGrief/AlliumLinter/src/config"
	"github.com/OnionGrief/AlliumLinter/src/lexer"
	"github.com/OnionGrief/AlliumLinter/src/logger"
	"unicode"
)

// CheckPrecondition Проверки которые возможно выполнить до построения дерева
func CheckPrecondition(tokens []lexer.Token) []logger.Log {
	var logs []logger.Log
	logs = append(logs, checkExtern(tokens)...)
	if config.Cfg.SnakeCase || config.Cfg.CamelCase {
		logs = append(logs, checkNames(tokens)...)
	}
	logs = append(logs, checkConstCount(tokens)...)
	return logs
}

// Проверка на вынесение констант
func checkConstCount(tokens []lexer.Token) []logger.Log {
	var logs []logger.Log
	count := config.Cfg.ConstCount
	if count < 3 {
		count = 3
	}
	leng := config.Cfg.ConstLen
	if leng < 3 {
		leng = 3
	}
	chks := make(map[string][]lexer.Position)
	str := ""
	var pos lexer.Position
	for _, tkn := range tokens {
		switch tkn.TokenType {
		case lexer.ASCIIV:
			if str == "" {
				pos = tkn.Start
			}
			str += tkn.Value
		default:
			if str != "" {
				if len(str) > int(leng) {
					chks[str] = append(chks[str], pos)
				}
				str = ""
			}
		}
	}
	for key, val := range chks {
		if len(val) > int(count) {
			logs = append(logs, logger.ExternLog(fmt.Sprintf("%s must be a const %s", key, prepareCoord(val))))
		}
	}
	return logs
}

func prepareCoord(in []lexer.Position) string {
	str := ""
	for _, val := range in {
		str = fmt.Sprintf("%s (%d:%d)", str, val.Row, val.Column)
	}
	return str
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
		if config.Cfg.CamelCase {
			if !isCamelCase(value) {
				logs = append(logs, logger.FormatNameLogCamel(tkn))
			}
		}
		if config.Cfg.SnakeCase {
			if !isSnakeCase(value) {
				logs = append(logs, logger.FormatNameLogShake(tkn))
			}
		}
	}
	return logs
}
func checkExtern(tokens []lexer.Token) []logger.Log {
	var logs []logger.Log
	var pos []lexer.Position
	for _, tkn := range tokens {
		if tkn.TokenType == lexer.EXTERN {
			pos = append(pos, tkn.Start)
		}
	}
	if len(pos) > 1 {
		logs = append(logs, logger.ExternLog(fmt.Sprintf("Обнаружено несколько EXTERN, рекомендуется использовать 1 %s", prepareCoord(pos))))
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
