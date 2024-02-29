package lexer

/*   legal tokens are   :                                       */
/*      VAR     -       variable,                               */
/*      ID      -       identifier,                             */
/*      VAR_type -   variable typed by a lower unicode-index    */
/*      PREWORD -       preword,                                */
/*      COMPSYM -       compound symbol,                        */
/*      MDIGIT  -       macrodigit,                             */
/*	    MSTRING -		  string not in ASCII Value			    */
/*      ASCIIV  -       symbol represented by its ASCII Value,  */
/*      LCBRAK  -       left concretization bracket,            */
/*      STRING  -       string of printable characters.         */
/*      EXTRN   -       ID of the form "EXTRN" or "ENTRY"       */
/*   In addition to these there is a number of other characters */
/*   which could be returned (list separated by blanks):        */
/*   = ( ) > , : ; { } [ ]                                      */
type TokenType uint

type Position struct {
	Column uint
	Row    uint
}
type Token struct {
	TokenType TokenType
	Value     string
	Start     Position
	Finish    Position
	Pos       uint64
}

func NewToken(val rune, start, finish Position) Token {
	return Token{TokenType: TokenType(val), Start: start, Finish: finish}
}
func NewTokenByType(val TokenType, start, finish Position) Token {
	return Token{TokenType: val, Start: start, Finish: finish}
}
func NewTokenByTypeAndVal(val TokenType, data string, start, finish Position) Token {
	return Token{TokenType: val, Value: data, Start: start, Finish: finish}
}

const (
	EOF        TokenType = iota //EOF
	NAME                        //имя?
	STRING                      //componoid
	MDIGIT                      //macrodogit
	ASCIIV                      //1 символ
	VAR                         //переменная
	LBRAK                       //(
	RBRAK                       //), s.Sort :
	LCBRAK                      //<
	RCBRAK                      //>
	OPENBLK                     //{
	CLOSEBLK                    //}
	COMMA                       //,
	COLON                       //:
	EQUAL                       //=
	SEMICOLON                   //;
	COMMENT                     //COMMENT *$
	BIGCOMMENT                  //COMMENT /* */
	EXTERN
	ENTRY
	SPACE
	NEWLINE
	TAB
)

func (tkn TokenType) String() string {
	switch tkn {
	case EOF:
		return "EOF"
	case NAME:
		return "NAME"
	case STRING:
		return "STRING"
	case MDIGIT:
		return "MDIGIT"
	case ASCIIV:
		return "ASCIIV"
	case VAR:
		return "VAR"
	case LBRAK:
		return "LBRAK"
	case RBRAK:
		return "RBRAK"
	case LCBRAK:
		return "LCBRAK"
	case RCBRAK:
		return "RCBRAK"
	case OPENBLK:
		return "OPENBLK"
	case CLOSEBLK:
		return "CLOSEBLK"
	case COMMA:
		return "COMMA"
	case COLON:
		return "COLON"
	case EQUAL:
		return "EQUAL"
	case SEMICOLON:
		return "SEMICOLON"
	case COMMENT:
		return "COMMENT"
	case BIGCOMMENT:
		return "BIGCOMMENT"
	case EXTERN:
		return "EXTERN"
	case ENTRY:
		return "ENTRY"
	case SPACE:
		return "SPACE"
	case NEWLINE:
		return "NEWLINE"
	case TAB:
		return "TAB"
	}
	return "ERROR"
}
