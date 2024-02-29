package tree

import (
	"fmt"
)

func CreateDot(root *treeElem) string {
	str := "digraph G {\n"
	str += dotFromElem(root, 0)
	str += "}\n"
	return str
}

var pos int = 0

func dotFromElem(root *treeElem, idx int) string {
	if root == nil {
		return fmt.Sprintf("%d [shape=Square,label=\"\"];\n", pos)
	}
	res := fmt.Sprintf("%d [shape=Square,label=\"", pos)
	res += string(root.typeEl) + "\n"
	curPos := pos
	if root.tkn != nil {
		res += root.tkn.Value + "\n"
		res += root.tkn.TokenType.String() + "\n"
		res += fmt.Sprintf("%d:%d - %d-%d\n", root.tkn.Start.Column, root.tkn.Start.Row, root.tkn.Finish.Column, root.tkn.Finish.Row)
	}
	res += "\"];\n"
	for _, val := range root.child {

		pos++
		ind := pos
		res += dotFromElem(val, ind)
		res += fmt.Sprintf("%d -> %d\n", curPos, ind)

	}
	return res
}
