package lexer

import (
	"bufio"
	"io"
)

type tknInfo struct {
	lastSymbol   symbol //храним последний символ
	needPrevious bool   //нужно ли вернуть прошлый символ
	Pos          Position
	lastSymboll  symbol
}

type lexer struct {
	reader *bufio.Reader
	info   *tknInfo
}

func NewLexer(reader *bufio.Reader) lexer {
	inf := &tknInfo{
		lastSymbol:   symbol{},
		needPrevious: false,
		Pos: Position{
			Column: 1,
			Row:    1,
		},
	}
	return lexer{reader: reader, info: inf}
}

type symbol struct {
	symbol rune
	EOF    bool // Для удобства обработки EOF
	pos    Position
}

func (l lexer) NewSymbol(symb rune, EOF bool) symbol {
	if EOF {
		return symbol{EOF: true, pos: l.info.Pos}
	}
	defer func() {
		if symb == '\n' {
			l.info.Pos.Row = 0
			l.info.Pos.Column++
		} else if symb == '\t' {
			l.info.Pos.Row += 2
		} else {
			l.info.Pos.Row++
		}
	}()
	pos1 := l.info.Pos
	return symbol{symbol: symb, pos: pos1}
}
func (l lexer) readChar() symbol {
	if l.info.needPrevious {
		l.info.needPrevious = false
		return l.info.lastSymbol
	}
	ch, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			l.info.lastSymbol = l.NewSymbol(0, true)
			return l.info.lastSymboll
		}
		panic(err)
	}
	l.info.lastSymbol = l.NewSymbol(ch, false)
	return l.info.lastSymbol
}

func (l lexer) GetToken() (res []Token, err error) {
	symb := l.readChar()
	if symb.symbol == 13 {
		return nil, nil
	}
	if symb.EOF {
		return []Token{NewTokenByType(EOF, symb.pos, symb.pos)}, nil
	}
	switch symb.symbol {
	case '(':
		return []Token{NewTokenByType(LBRAK, symb.pos, symb.pos)}, nil
	case '\t':
		return []Token{NewTokenByType(TAB, symb.pos, symb.pos)}, nil
	case ')':
		return []Token{NewTokenByType(RBRAK, symb.pos, symb.pos)}, nil
	case '<':
		return []Token{NewTokenByType(LCBRAK, symb.pos, symb.pos)}, nil
	case '>':
		return []Token{NewTokenByType(RCBRAK, symb.pos, symb.pos)}, nil
	case '{':
		return []Token{NewTokenByType(OPENBLK, symb.pos, symb.pos)}, nil
	case '}':
		return []Token{NewTokenByType(CLOSEBLK, symb.pos, symb.pos)}, nil
	case ',':
		return []Token{NewTokenByType(COMMA, symb.pos, symb.pos)}, nil
	case ':':
		return []Token{NewTokenByType(COLON, symb.pos, symb.pos)}, nil
	case '=':
		return []Token{NewTokenByType(EQUAL, symb.pos, symb.pos)}, nil
	case ';':
		return []Token{NewTokenByType(SEMICOLON, symb.pos, symb.pos)}, nil
	case '\n':
		return []Token{NewTokenByType(NEWLINE, symb.pos, symb.pos)}, nil
	case ' ':
		return []Token{NewTokenByType(SPACE, symb.pos, symb.pos)}, nil
	case '$':
		symb1 := l.readChar()
		if symb1.EOF {
			return []Token{}, ErrInLexer
		}
		if symb1.symbol == 'E' {
			symb2 := l.readChar()
			if symb2.EOF {
				return []Token{}, ErrInLexer
			}
			switch symb2.symbol {
			case 'X':
				checkStr := "TERN"
				symb3 := symb2
				for len(checkStr) > 0 {
					symb3 = l.readChar()
					if symb3.EOF {
						return []Token{}, ErrInLexer
					}
					if symb3.symbol != rune(checkStr[0]) {
						return []Token{}, ErrInLexer
					}
					checkStr = checkStr[1:]
				}
				symb4 := l.readChar()
				if symb4.EOF {
					return []Token{NewTokenByType(EXTERN, symb.pos, symb3.pos), NewTokenByType(EOF, symb4.pos, symb4.pos)}, nil
				} else if isAlnum(symb4.symbol) {
					return []Token{}, ErrInLexer
				} else {
					l.info.needPrevious = true
					return []Token{NewTokenByType(EXTERN, symb.pos, symb3.pos)}, nil
				}
			case 'N':
				checkStr := "TRY"
				symb3 := symb2
				for len(checkStr) > 0 {
					symb3 = l.readChar()
					if symb3.EOF {
						return []Token{}, ErrInLexer
					}
					if symb3.symbol != rune(checkStr[0]) {
						return []Token{}, ErrInLexer
					}
					checkStr = checkStr[1:]
				}
				symb4 := l.readChar()
				if symb4.EOF {
					return []Token{NewTokenByType(ENTRY, symb.pos, symb3.pos), NewTokenByType(EOF, symb4.pos, symb4.pos)}, nil
				} else if isAlnum(symb4.symbol) {
					return []Token{}, ErrInLexer
				} else {
					l.info.needPrevious = true
					return []Token{NewTokenByType(ENTRY, symb.pos, symb3.pos)}, nil
				}
			default:
				return []Token{}, ErrInLexer
			}
		} else {
			return []Token{}, ErrInLexer
		}
	case '/':
		symb1 := l.readChar()
		if symb1.EOF {
			return []Token{}, ErrInLexer
		}
		if symb1.symbol == '*' {
			comment := ""
			for {
				symb2 := l.readChar()
				if symb2.EOF {
					return []Token{}, ErrInLexer
				}
				if symb2.symbol == '*' {
					symb3 := l.readChar()
					if symb3.EOF {
						return []Token{}, ErrInLexer
					}
					if symb3.symbol == '/' {
						return []Token{NewTokenByTypeAndVal(BIGCOMMENT, comment, symb.pos, symb3.pos)}, nil
					} else {
						l.info.needPrevious = true
					}
				} else {
					comment += string(symb2.symbol)
				}
			}
		} else {
			return []Token{}, ErrInLexer
		}
	case '*':
		comment := ""
		for {
			pos1 := symb.pos
			symb2 := l.readChar()
			if symb2.EOF {
				return []Token{}, ErrInLexer
			}
			if symb2.symbol == '\n' {
				l.info.needPrevious = true
				return []Token{NewTokenByTypeAndVal(COMMENT, comment, symb.pos, pos1)}, nil
			} else {
				comment += string(symb2.symbol)
			}
		}

	case 's', 't', 'e':
		symb1 := l.readChar()
		if symb1.EOF {
			return []Token{}, ErrInLexer
		}
		symb2 := symb1
		if symb1.symbol == '.' {
			//Переменная
			str := string(symb.symbol) + "."
			for {
				pos1 := symb2.pos
				symb2 = l.readChar()
				if symb2.EOF {
					return []Token{}, ErrInLexer
				}
				if !isAlnum(symb2.symbol) {
					if len(str) > 0 {
						l.info.needPrevious = true
						return []Token{NewTokenByTypeAndVal(VAR, str, symb.pos, pos1)}, nil
					} else {
						return []Token{}, ErrInLexer
					}
				} else {
					str += string(symb2.symbol)
				}
			}
		} else {
			str := string(symb.symbol)
			symb2 := symb1
			for {
				pos1 := symb2.pos
				symb2 = l.readChar()
				if symb2.EOF {
					return []Token{}, ErrInLexer
				}
				if !isAlnum(symb2.symbol) {
					l.info.needPrevious = true
					return []Token{NewTokenByTypeAndVal(NAME, str, symb.pos, pos1)}, nil
				} else {
					str += string(symb2.symbol)
				}
			}
		}
	case '\'':
		var tkns []Token
		for {
			symb1 := l.readChar()
			if symb1.EOF {
				return []Token{}, ErrInLexer
			}
			//(c == 'n' || c == 'r' || c == 't' || c == '\'' || c == '"' || ' ' ||
			//				c == '<' || c == '>' || c == '(' || c == ')' || c == ESCAPE)
			if symb1.symbol == '\\' {
				symb2 := l.readChar()
				if symb2.EOF {
					return []Token{}, ErrInLexer
				}
				switch symb2.symbol {
				case '\'', 'n', 'r', 't', '"', ' ', '<', '>', '(', ')', '\\':
					res := "\\" + string(symb2.symbol)
					tkns = append(tkns, NewTokenByTypeAndVal(ASCIIV, res, symb.pos, symb1.pos))
				default:
					return []Token{}, ErrInLexer
				}
			} else if symb1.symbol == '\'' {
				return tkns, nil
			} else if !isFormatSymb(symb1.symbol) && !isAlnum(symb1.symbol) {
				return []Token{}, ErrInLexer
			} else {
				tkns = append(tkns, NewTokenByTypeAndVal(ASCIIV, string(symb1.symbol), symb.pos, symb1.pos))
			}
		}
	case '"':
		res := ""
		for {
			symb1 := l.readChar()
			if symb1.EOF {
				return []Token{}, ErrInLexer
			}
			if symb1.symbol == '"' {
				return []Token{NewTokenByTypeAndVal(STRING, res, symb.pos, symb1.pos)}, nil
			} else if !isAlnum(symb1.symbol) {
				return []Token{}, ErrInLexer
			} else {
				res += string(symb1.symbol)
			}
		}
	default:
		str := ""
		if isDigit(symb.symbol) {
			str += string(symb.symbol)
			symb1 := symb
			for {
				pos1 := symb1.pos
				symb1 := l.readChar()
				if symb1.EOF {
					return []Token{}, ErrInLexer
				}
				if isDigit(symb1.symbol) {
					str += string(symb1.symbol)
				} else {
					l.info.needPrevious = true
					return []Token{NewTokenByTypeAndVal(MDIGIT, str, symb.pos, pos1)}, nil
				}
			}
		}
		if isAlpha(symb.symbol) {
			str += string(symb.symbol)
			symb1 := symb
			for {
				pos1 := symb1.pos
				symb1 = l.readChar()
				if symb1.EOF {
					return []Token{}, ErrInLexer
				}
				if isAlnum(symb1.symbol) {
					str += string(symb1.symbol)
				} else {
					l.info.needPrevious = true
					return []Token{NewTokenByTypeAndVal(NAME, str, symb.pos, pos1)}, nil
				}
			}
		}
	}
	return []Token{NewTokenByType(EOF, symb.pos, symb.pos)}, nil
}

func (l lexer) GetAllTokens() (res []Token, err error) {
	defer func() {
		if err == nil {
			for i, _ := range res {
				res[i].Pos = uint64(i)
			}
		}
	}()
	var tkns []Token
	for {
		l, err := l.GetToken()
		if err != nil {
			return nil, err
		}
		tkns = append(tkns, l...)
		if len(tkns) > 0 && tkns[len(tkns)-1].TokenType == EOF {
			return tkns, nil
		}
	}
}

func DeleteSugar(in []Token) []Token {
	var tkns []Token
	for _, tkn := range in {
		if tkn.TokenType != SPACE && tkn.TokenType != NEWLINE && tkn.TokenType != TAB && tkn.TokenType != BIGCOMMENT && tkn.TokenType != COMMENT {
			tkns = append(tkns, tkn)
		}
	}
	return tkns
}
