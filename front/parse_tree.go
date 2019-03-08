package front

import "fmt"

type ParseTree struct {
	Nodes []ParseTreeNode
}

func (p ParseTree) String() string {
	var result string
	for _, node := range p.Nodes {
		result += fmt.Sprintf("%s\n", node.Print())
	}
	return result
}

type ParseTreeNode interface {
	NodeName() string
	Print() string
}
