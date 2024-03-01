package tree

import (
	"github.com/OnionGrief/AlliumLinter/src/config"
	"github.com/OnionGrief/AlliumLinter/src/lexer"
	"github.com/OnionGrief/AlliumLinter/src/logger"
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
	for _, key := range config.UsingFunctions {
		p.usingFunctions[key] = struct{}{}
	}
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
		if i == 0 {
			return true
		}
		i--
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
			substrings := []string{"s.", "e.", "t."}
			totalCount := 0
			for _, substr := range substrings {
				count := strings.Count(val.Value, substr)
				totalCount += count
			}
			if totalCount > 3 {

				res = append(res, logger.CodeInComment(val))
			}
		}
	}
	return res
}

type сoords struct {
	Start lexer.Position
	End   lexer.Position
}

func isCoordsNested(coords сoords, otherCoords сoords) bool {
	return coords.Start.Column >= otherCoords.Start.Column && coords.Start.Row >= otherCoords.Start.Row &&
		coords.End.Column <= otherCoords.End.Column && coords.End.Row <= otherCoords.End.Row
}

type forCheckTreeRec struct {
	first, second сoords
	deep          int
}

func findNonNestedCoords(in []forCheckTreeRec) []logger.Log {
	var logs []logger.Log
	for i, coords := range in {
		isNested := false
		for j, otherCoords := range in {
			if i != j && isCoordsNested(coords.first, otherCoords.first) {
				isNested = true
				break
			}
		}
		if !isNested {
			logs = append(logs, logger.ReusingBlock(coords.deep, coords.first.Start, coords.second.Start))
		}
	}
	return logs
}

func CheckTreeRec(elem *treeElem) []logger.Log {
	var forCheck []forCheckTreeRec
	collection, _ := CollectTreeRec(elem)
	if len(collection) == 0 {
		return nil
	}
	combinations := generateUniqueCombinations(len(collection) - 1)

	for _, comb := range combinations {
		if r := checkTreeRecBool(collection[comb.a].elem, collection[comb.b].elem); r {
			if collection[comb.b].deep > int(config.BlockLen) {
				forCheck = append(forCheck, forCheckTreeRec{
					first: сoords{
						Start: collection[comb.a].elem.start,
						End:   collection[comb.a].elem.finish,
					},
					second: сoords{
						Start: collection[comb.b].elem.start,
						End:   collection[comb.b].elem.finish,
					},
					deep: collection[comb.a].deep,
				})
			}
		}
	}

	return findNonNestedCoords(forCheck)
}

type checkTree struct {
	elem *treeElem
	deep int
}

func CollectTreeRec(elem *treeElem) ([]checkTree, int) {
	var out []checkTree
	i := 0
	for _, child := range elem.child {
		if child != nil {
			nodes, deep := CollectTreeRec(child)
			out = append(out, nodes...)
			if deep > i {
				i = deep
			}

		}
	}
	if elem.typeEl != PatternExprTerm {
		out = append(out, checkTree{
			elem: elem,
			deep: i + 1,
		})
	}
	return out, i + 1
}

func checkTreeRecBool(elem1, elem2 *treeElem) bool {
	if elem1 == nil || elem2 == nil {
		return false
	}
	if !((elem1.tkn == nil && elem2.tkn == nil) ||
		(elem1.tkn != nil && elem2.tkn != nil &&
			elem1.tkn.TokenType == elem2.tkn.TokenType && elem1.tkn.Value == elem2.tkn.Value)) {
		return false
	}
	if len(elem1.child) != len(elem2.child) {
		return false
	}
	for idx, _ := range elem1.child {
		if !checkTreeRecBool(elem1.child[idx], elem2.child[idx]) {
			return false
		}
	}
	return true
}

type combination struct {
	a, b int
}

func generateUniqueCombinations(n int) []combination {
	var out []combination
	for i := 0; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			out = append(out, combination{
				a: i,
				b: j,
			})
		}
	}
	return out
}
