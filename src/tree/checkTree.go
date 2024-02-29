package tree

import (
	"refalLint/src/config"
	"refalLint/src/lexer"
	"refalLint/src/logger"
	"strings"
)

func (p *Parser) Checks() []logger.Log {
	var res []logger.Log
	res = append(res, p.CheckComment()...)
	res = append(res, p.CheckUsingFunctions()...)
	return res
}

func (p *Parser) CheckUsingFunctions() []logger.Log {
	var res []logger.Log
	for key, idx := range p.declaredFunctions {
		if _, exists := p.usingFunctions[key]; !exists {
			if !p.searchComments(idx) {
				res = append(res, logger.UnreachableLog(p.allTokens[idx]))
			}
		}
	}
	return res
}

func (p *Parser) searchComments(idx uint64) bool {
	i := idx
	for {
		i--
		if i < 0 {
			return true
		}
		tkn := p.allTokens[i]
		if tkn.TokenType == lexer.COMMENT {
			if strings.Contains(tkn.Value, config.FuncComment) {
				return true
			}
			return false
		}
		if tkn.TokenType != lexer.SPACE && tkn.TokenType != lexer.NEWLINE && tkn.TokenType != lexer.TAB {
			return false
		}
	}
}

type funcCoords struct {
	start, finish uint64
}

func isNumberInRange(number uint64, coord funcCoords) bool {
	return number >= coord.start && number <= coord.finish
}

func (p *Parser) CheckComment() []logger.Log {
	var res []logger.Log
	for _, val := range p.allTokens {
		if val.TokenType == lexer.BIGCOMMENT {
			for _, coords := range p.funcCoords {
				if isNumberInRange(val.Pos, coords) {
					res = append(res, logger.CodeInComment(val))
				}
			}
		}
	}
	return res
}
