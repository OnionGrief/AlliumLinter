package lexer

import (
	"unicode"
)

func isAlnum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' || r == '@'
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '-' || r == '_' || r == '.' || r == '@'
}
func isFormatSymb(r rune) bool {
	return r != '\''
	//return r == ' ' || r == ':' || r == '=' || r == '(' || r == ')' || r == '+' || r == '-' || r == '/' || r == '\\' || r == '{' || r == '}' || r == '<' || r == '>' || r == ','
}

func isCntrl(r rune) bool {
	return unicode.IsControl(r)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func isLower(r rune) bool {
	return unicode.IsLower(r)
}

func isGraph(r rune) bool {
	return unicode.IsGraphic(r)
}

func isPrint(r rune) bool {
	return unicode.IsPrint(r)
}

func isPunct(r rune) bool {
	return unicode.IsPunct(r)
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}

func isUpper(r rune) bool {
	return unicode.IsUpper(r)
}

func isXDigit(r rune) bool {
	return unicode.Is(unicode.Hex_Digit, r)
}
func xDigitToInt(r rune) int32 {
	res := int32(0)
	r = unicode.ToLower(r)
	if r <= '9' {
		res = r - '0'
	} else {
		res = r - 'a' + 10
	}
	return res
}
