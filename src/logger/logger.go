package logger

import (
	"fmt"
	"refalLint/src/lexer"
	"strings"
)

type logType uint
type LogLevel uint

const (
	formatName logType = iota
	extern

	Warning LogLevel = iota
)

type Log struct {
	str    string
	lgType logType
	level  LogLevel
}

func (l Log) String() string {
	return l.str
}

func FormatNameLogCamel(tkn lexer.Token) Log {
	return Log{
		str:    fmt.Sprintf("%s must be camelCase pos(%d:%d-%d%d)", tkn.Value, tkn.Start.Column, tkn.Start.Row, tkn.Finish.Column, tkn.Finish.Row),
		lgType: formatName,
		level:  Warning,
	}
}

func FormatNameLogShake(tkn lexer.Token) Log {
	return Log{
		str:    fmt.Sprintf("%s must be snake_case pos(%d:%d-%d:%d)", tkn.Value, tkn.Start.Column, tkn.Start.Row, tkn.Finish.Column, tkn.Finish.Row),
		lgType: formatName,
		level:  Warning,
	}
}
func ExternLog(str string) Log {
	return Log{
		str:    str,
		lgType: extern,
		level:  Warning,
	}
}
func UnreachableLog(tkn lexer.Token) Log {
	return Log{
		str:    fmt.Sprintf("Недостижимая функция %s (%d:%d-%d:%d)", tkn.Value, tkn.Start.Column, tkn.Start.Row, tkn.Finish.Column, tkn.Finish.Row),
		lgType: extern,
		level:  Warning,
	}
}
func CodeInComment(tkn lexer.Token) Log {
	if strings.Contains(tkn.Value, "\n") {
		strs := strings.Split(tkn.Value, "\n")
		if strs[0] == "" {
			strs = strs[1:]
		}
		return Log{
			str:    fmt.Sprintf("Код в комментарии %s (%d:%d-%d:%d)", strs[0], tkn.Start.Column, tkn.Start.Row, tkn.Finish.Column, tkn.Finish.Row),
			lgType: extern,
			level:  Warning,
		}
	}
	return Log{
		str:    fmt.Sprintf("Код в комментарии %s (%d:%d-%d:%d)", tkn.Value, tkn.Start.Column, tkn.Start.Row, tkn.Finish.Column, tkn.Finish.Row),
		lgType: extern,
		level:  Warning,
	}

}
func ReusingBlock(len int, pos1, pos2 lexer.Position) Log {
	return Log{
		str:    fmt.Sprintf("Найден переиспользуемый блок длины %d [%d:%d] [%d:%d]", len, pos1.Column, pos1.Row, pos2.Column, pos2.Row),
		lgType: extern,
		level:  Warning,
	}
}
