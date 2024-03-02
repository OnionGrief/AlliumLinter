package tree

type ElemType string

var (
	Program           ElemType = "Program"
	ExternFuncs       ElemType = "ExternFuncs"
	Functions         ElemType = "Functions"
	ExternFunc        ElemType = "ExternFunc"
	EXTERN            ElemType = "EXTERN"
	FunctionName      ElemType = "FunctionName"
	Function          ElemType = "Function"
	FunctionBody      ElemType = "FunctionBody"
	ENTRY             ElemType = "ENTRY"
	IDENT             ElemType = "IDENT"
	Sentences         ElemType = "Sentences"
	Sentence          ElemType = "Sentence"
	PatternExpr       ElemType = "PatternExpr"
	Conditions        ElemType = "Conditions"
	RightSentencePart ElemType = "RightSentencePart"
	ResultExpr        ElemType = "ResultExpr"
	ResultExprTerm    ElemType = "ResultExprTerm"
	PatternExprTerm   ElemType = "PatternExprTerm"
	STRING            ElemType = "STRING"
	INTEGER           ElemType = "INTEGER"
	VARNAME           ElemType = "VARNAME"
	SPECIAL           ElemType = "SPECIAL"
)
