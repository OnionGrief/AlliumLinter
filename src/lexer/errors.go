package lexer

import "fmt"

type LexerError error

var (
	ErrInLexer               LexerError = fmt.Errorf("ErrInLexer")
	ErrInParseEntryFunc                 = fmt.Errorf("error in parseEntryFunc")
	ErrInParseFunc                      = fmt.Errorf("error in parseFunc")
	ErrInParseSentence                  = fmt.Errorf("error in parseSentence")
	ErrInParseSentenceRest              = fmt.Errorf("error in parseSentenceRest")
	ErrInParseConditionsRest            = fmt.Errorf("error in parseConditionsRest")
	ErrInParseExternFunc                = fmt.Errorf("error in parseExternFunc")
	ErrInParseResultExprTerm            = fmt.Errorf("error in parseResultExprTerm")
)
