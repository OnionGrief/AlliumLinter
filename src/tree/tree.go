package tree

import (
	"errors"
	"refalLint/src/lexer"
)

/*
Program → ExternFuncs Program

	| Functions

ExternFuncs -> ExternFuncs ExternFunc

	| ε

ExternFunc -> "$EXTERN" FunctionName ";"
Functions → Functions Function

	| ε

Function → FunctionName FunctionBody

	|  "$ENTRY" FunctionName FunctionBody

FunctionName → IDENT
FunctionBody → '{' Sentences '}'

Sentences → Sentences Sentence

	| ε

Sentence            → PatternExpr ',' SentenceRest ';'
SentenceRest        → ResultExpr ':' '{' Sentences '}'

	| ResultExpr ':' PatternExpr ConditionsRest
	| RightSentencePartWithEqual

ConditionsRest      → ',' ResultExpr ':' PatternExpr ConditionsRest

	| RightSentencePartWithEqual

RightSentencePartWithEqual → '=' ResultExpr
Conditions → Conditions ',' ResultExpr ':' PatternExpr

	| ε

PatternExpr → PatternExpr PatternExprTerm

	| ε

PatternExprTerm → STRING

	| INTEGER
	| IDENT
	| VARNAME
	| '(' PatternExpr ')'

ResultExpr → ResultExpr ResultExprTerm

	| ε

ResultExprTerm → STRING

	| INTEGER
	| IDENT
	| VARNAME
	| '(' ResultExpr ')'
	| '<' IDENT ResultExpr '>'
*/
type treeElem struct {
	typeEl ElemType
	tkn    *lexer.Token
	child  []*treeElem
}

type Parser struct {
	tokens            []lexer.Token       // Слайс токенов, полученных от лексера без мусора
	allTokens         []lexer.Token       // Слайс токенов, полученных от лексера
	current           int                 // Текущий индекс в слайсе токенов
	usingFunctions    map[string]struct{} //	Функции которые вызываются
	declaredFunctions map[string]uint64   //	Функции которые обьявлены (и номер токена)
	funcCoords        []funcCoords
}

// NewParser Функция для создания нового парсера
func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:            lexer.DeleteSugar(tokens),
		allTokens:         tokens,
		current:           0,
		usingFunctions:    make(map[string]struct{}),
		declaredFunctions: make(map[string]uint64),
	}
}
func (p *Parser) NewSpecial(tkn lexer.Token) *treeElem {
	return &treeElem{tkn: &tkn, typeEl: SPECIAL}
}
func NewWithType(tkn lexer.Token, elem ElemType) *treeElem {
	return &treeElem{tkn: &tkn, typeEl: elem}
}

func (p *Parser) Parse() (*treeElem, error) {
	if len(p.tokens) == 0 {
		return nil, errors.New("no tokens to parse")
	}
	ast := &treeElem{typeEl: "root"}
	// Основной цикл парсинга
	for !p.isAtEnd() {
		node, err := p.parseUnit()
		if err != nil {
			return nil, err
		}
		ast.child = append(ast.child, node) // Добавляем разобранный узел в AST
	}

	return ast, nil
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)-1
}
func (p *Parser) getCurTkn() lexer.Token {
	return p.tokens[p.current]
}
func (p *Parser) nexttkn() {
	p.current++
}
func (p *Parser) parseUnit() (*treeElem, error) {
	tkn := p.getCurTkn()
	switch tkn.TokenType {
	case lexer.BIGCOMMENT, lexer.COMMENT:
		p.nexttkn()
		return p.NewSpecial(tkn), nil
	case lexer.EXTERN:
		p.nexttkn()
		return p.parseExternFunc()
	case lexer.ENTRY:
		return p.parseEntryFunc()
	case lexer.NAME:
		return p.parseFunc()
	default:
		return nil, nil

	}
	return nil, nil
}
func (p *Parser) parseExternFunc() (*treeElem, error) {
	ast := &treeElem{typeEl: ExternFunc}
	tkn := p.getCurTkn()
	ast.tkn = &tkn
	p.nexttkn()
	for {
		tkn1 := p.getCurTkn()
		p.nexttkn()
		switch tkn1.TokenType {
		case lexer.NAME, lexer.COMMA:
			ast.child = append(ast.child, NewWithType(tkn1, ExternFunc))
		case lexer.SEMICOLON:
			ast.child = append(ast.child, NewWithType(tkn1, ExternFunc))
			return ast, nil
		default:
			return nil, lexer.ErrInParseExternFunc
		}
	}
	return nil, lexer.ErrInParseExternFunc
}
func (p *Parser) parseEntryFunc() (*treeElem, error) {
	ast := &treeElem{typeEl: ENTRY}
	tkn := p.getCurTkn()
	ast.tkn = &tkn
	p.nexttkn()
	for {
		tkn1 := p.getCurTkn()
		switch tkn1.TokenType {
		case lexer.NAME:
			node, err := p.parseFunc()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, node)

			return ast, nil
		default:
			return nil, lexer.ErrInParseEntryFunc
		}
	}
	return nil, lexer.ErrInParseEntryFunc
}
func (p *Parser) parseFunc() (*treeElem, error) {
	ast := &treeElem{typeEl: "ParseFunc"}

	tkn := p.getCurTkn()
	ast.tkn = &tkn
	p.funcCoords = append(p.funcCoords, funcCoords{start: tkn.Pos})
	p.nexttkn()
	switch tkn.TokenType {
	case lexer.NAME:
		p.declaredFunctions[tkn.Value] = tkn.Pos //Пишем в используемые функции
		for {
			tkn1 := p.getCurTkn()
			p.nexttkn()
			switch tkn1.TokenType {
			case lexer.OPENBLK:
				ast.child = append(ast.child, NewWithType(tkn1, "ParseFunc"))
				for {
					tkn2 := p.getCurTkn()
					switch tkn2.TokenType {
					case lexer.CLOSEBLK:
						p.nexttkn()
						ast.child = append(ast.child, NewWithType(tkn2, "ParseFunc"))
						tkn2 := p.getCurTkn()
						p.funcCoords[len(p.funcCoords)-1].finish = tkn2.Pos - 1
						return ast, nil
					default:
						node, err := p.parseFuncCore()
						if err != nil {
							return nil, err
						}
						ast.child = append(ast.child, node)
					}
				}
			default:
				return nil, lexer.ErrInParseFunc
			}
		}
	default:
		return nil, lexer.ErrInParseFunc
	}
	return ast, nil
}

func (p *Parser) parseFuncCore() (*treeElem, error) {
	node, err := p.parseSentence()
	if err != nil {
		return nil, err
	}
	return node, nil
}
func (p *Parser) parseSentence() (*treeElem, error) {
	ast := &treeElem{typeEl: Sentence}
	node, err := p.parsePatternExpr()
	if err != nil {
		return nil, err
	}
	ast.child = append(ast.child, node)

	tkn2 := p.getCurTkn()
	if tkn2.TokenType != lexer.COMMA && tkn2.TokenType != lexer.EQUAL {
		return nil, lexer.ErrInParseSentence
	}

	node1, err := p.parseSentenceRest()
	if err != nil {
		return nil, err
	}
	ast.child = append(ast.child, node1)
	tkn1 := p.getCurTkn()
	p.nexttkn()
	if tkn1.TokenType != lexer.SEMICOLON {
		return nil, lexer.ErrInParseSentence
	}
	ast.child = append(ast.child, NewWithType(tkn1, Sentence))
	return ast, nil
}

func (p *Parser) parseSentenceRest() (*treeElem, error) {
	ast := &treeElem{typeEl: Conditions}
	for {
		tkn := p.getCurTkn()
		switch tkn.TokenType {
		case lexer.SEMICOLON:
			{
				return ast, nil
			}
		case lexer.EQUAL:
			{
				elem, err := p.parseRightSentencePartWithEqual()
				if err != nil {
					return nil, err
				}
				ast.child = append(ast.child, elem)
				return ast, nil
			}
		default:
			ast.child = append(ast.child, NewWithType(tkn, Conditions))
			p.nexttkn()
			tkn, err := p.parseResultExpr()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, tkn)
			tkn1 := p.getCurTkn()
			if tkn1.TokenType != lexer.COLON {
				return nil, lexer.ErrInParseSentenceRest
			}
			p.nexttkn()
			tkn2 := p.getCurTkn()

			if tkn2.TokenType == lexer.OPENBLK {
				p.nexttkn()
				for {
					tkn3 := p.getCurTkn()
					switch tkn3.TokenType {
					case lexer.CLOSEBLK:
						p.nexttkn()
						ast.child = append(ast.child, NewWithType(tkn3, "ParseFunc"))
						return ast, nil
					default:
						node, err := p.parseFuncCore()
						if err != nil {
							return nil, err
						}
						ast.child = append(ast.child, node)
					}
				}
			} else {
				node, err := p.parsePatternExpr()
				if err != nil {
					return nil, err
				}
				ast.child = append(ast.child, node)
				node1, err := p.parseConditionsRest()
				if err != nil {
					return nil, err
				}
				ast.child = append(ast.child, node1)
			}

		}
	}
}
func (p *Parser) parseConditionsRest() (*treeElem, error) {
	ast := &treeElem{typeEl: Sentences}
	for {
		tkn := p.getCurTkn()
		switch tkn.TokenType {

		case lexer.COMMA:
			ast.child = append(ast.child, NewWithType(tkn, RightSentencePart))
			p.nexttkn()
			tkn1, err := p.parseResultExpr()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, tkn1)
			tkn2 := p.getCurTkn()
			if tkn2.TokenType != lexer.COLON {
				return nil, lexer.ErrInParseConditionsRest
			}
			p.nexttkn()
			tkn21 := p.getCurTkn()
			if tkn21.TokenType == lexer.OPENBLK {
				for {
					tkn3 := p.getCurTkn()
					switch tkn3.TokenType {
					case lexer.CLOSEBLK:
						p.nexttkn()
						ast.child = append(ast.child, NewWithType(tkn3, "ParseFunc"))
						return ast, nil
					default:
						node, err := p.parseFuncCore()
						if err != nil {
							return nil, err
						}
						ast.child = append(ast.child, node)
					}
				}
			} else {
				node, err := p.parsePatternExpr()
				if err != nil {
					return nil, err
				}
				ast.child = append(ast.child, node)
				node1, err := p.parseConditionsRest()
				if err != nil {
					return nil, err
				}
				ast.child = append(ast.child, node1)
			}
			return ast, nil
		default:
			return p.parseRightSentencePartWithEqual()
		}
	}
}
func (p *Parser) parseRightSentencePartWithEqual() (*treeElem, error) {
	ast := &treeElem{typeEl: Sentences}
	for {
		tkn := p.getCurTkn()
		switch tkn.TokenType {

		case lexer.SEMICOLON:
			ast.child = append(ast.child, NewWithType(tkn, RightSentencePart))
			p.nexttkn()
			return ast, nil
		case lexer.EQUAL:
			ast.child = append(ast.child, NewWithType(tkn, RightSentencePart))
			p.nexttkn()
			tkn1, err := p.parseResultExpr()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, tkn1)
			return ast, nil
		default:
			return p.parseResultExpr()
		}
	}
}

func (p *Parser) parseResultExpr() (*treeElem, error) {
	ast := &treeElem{typeEl: ResultExpr}
	for {
		tkn := p.getCurTkn()
		switch tkn.TokenType {
		case lexer.RCBRAK, lexer.RBRAK:
			return ast, nil
		case lexer.SEMICOLON, lexer.COLON:
			return ast, nil
		default:
			node, err := p.parseResultExprTerm()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, node)
		}
	}
}

func (p *Parser) parseResultExprTerm() (*treeElem, error) {
	ast := &treeElem{typeEl: "ResultExprTerm"}
	for {
		tkn := p.getCurTkn()

		switch tkn.TokenType {
		case lexer.RBRAK:
			return nil, nil
		case lexer.STRING, lexer.MDIGIT, lexer.VAR, lexer.NAME, lexer.ASCIIV:
			ast.tkn = &tkn
			p.nexttkn()
			return ast, nil
		case lexer.LBRAK:
			p.nexttkn()
			ast.child = append(ast.child, NewWithType(tkn, "ResultExprTerm"))
			node, err := p.parseResultExpr()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, node)
			tkn1 := p.getCurTkn()

			p.nexttkn()
			if tkn1.TokenType == lexer.RBRAK {
				ast.child = append(ast.child, NewWithType(tkn1, "ResultExprTerm"))
				return ast, nil
			}
		case lexer.LCBRAK:
			p.nexttkn()
			ast.child = append(ast.child, NewWithType(tkn, "ResultExprTerm"))
			tkn1 := p.getCurTkn()

			if tkn1.TokenType != lexer.NAME {
				return nil, lexer.ErrInParseResultExprTerm
			}
			p.usingFunctions[tkn1.Value] = struct{}{} //Пишем функции которые вызвали
			ast.child = append(ast.child, NewWithType(tkn1, "ResultExprTerm"))
			p.nexttkn()
			node, err := p.parseResultExpr()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, node)
			tkn2 := p.getCurTkn()
			if tkn2.TokenType == lexer.RCBRAK {
				p.nexttkn()
				ast.child = append(ast.child, NewWithType(tkn2, "ResultExprTerm"))
				return ast, nil
			}
		default:
			return nil, nil
		}
	}
}

func (p *Parser) parsePatternExprTerm() (*treeElem, error) {
	ast := &treeElem{typeEl: PatternExprTerm}
	for {
		tkn := p.getCurTkn()
		p.nexttkn()
		switch tkn.TokenType {
		case lexer.STRING, lexer.MDIGIT, lexer.VAR, lexer.NAME, lexer.ASCIIV:
			ast.tkn = &tkn
			return ast, nil
		case lexer.LBRAK:
			ast.child = append(ast.child, NewWithType(tkn, PatternExprTerm))
			node, err := p.parsePatternExpr()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, node)
			tkn1 := p.getCurTkn()

			p.nexttkn()
			if tkn1.TokenType == lexer.RBRAK {
				ast.child = append(ast.child, NewWithType(tkn1, PatternExprTerm))
				return ast, nil
			}
		default:
			return nil, nil
		}
	}

}

func (p *Parser) parsePatternExpr() (*treeElem, error) {
	ast := &treeElem{typeEl: PatternExpr}
	for {
		tkn := p.getCurTkn()
		switch tkn.TokenType {
		case lexer.EQUAL, lexer.RBRAK:
			return ast, nil
		case lexer.COMMA:
			return ast, nil
		default:
			node, err := p.parsePatternExprTerm()
			if err != nil {
				return nil, err
			}
			ast.child = append(ast.child, node)
		}
	}

}
