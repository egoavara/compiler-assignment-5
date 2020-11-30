package cdtgo

import (
	"fmt"
	"strings"
)

type Node struct {
	data  interface{}
	next  *Node
	child *Node
}

func (s *Node) IsTerminal() bool {
	if _, ok := s.data.(Token); ok {
		return true
	} else {
		return false
	}
}
func (s *Node) String() string {
	if s.IsTerminal() {
		return fmt.Sprintf("terminal(%v)", s.data)
	} else {
		return fmt.Sprintf("nonterminal(%v)", s.data)
	}
}

func Parse(src string) (*Node, error) {
	tk, err := Scanning(src)
	if err != nil {
		return nil, err
	}
	var stack []*Node
	parser := NewParser(tk...)
	parser.ShiftHandle = func(token Token) {
		if isMean(token) {
			stack = append(stack, terminalNode(token))
		} else {
			stack = append(stack, nil)
		}
	}
	parser.ReduceHandle = func(rule Rule) {
		fold := filtered(stack[len(stack)-rightLength[rule]:])
		//
		stack = stack[:len(stack)-rightLength[rule]]
		for i := 0; i < len(fold)-1; i++ {
			fold[i].next = fold[i+1]
		}
		var first *Node = nil
		if len(fold) > 0 {
			first = fold[0]
		}
		if rule.IsNaming() {
			tmp := nonterminalNode(rule)
			tmp.child = first
			stack = append(stack, tmp)
		} else {
			stack = append(stack, first)
		}
	}
	err = parser.Parsing()
	if err != nil {
		return nil, err
	}
	return stack[0], nil
}

func terminalNode(tk Token) *Node {
	return &Node{
		data: tk,
	}
}
func nonterminalNode(rule Rule) *Node {

	return &Node{
		data: rule,
	}
}
func isMean(tk Token) bool {
	switch tk.Kind {
	case Tident:
		fallthrough
	case Tnumber:
		return true
	default:
		return false
	}
}
func filtered(nodes []*Node) []*Node {
	res := make([]*Node, 0, len(nodes))
	for _, node := range nodes {
		if node != nil {
			res = append(res, node)
		}
	}
	return res
}
func (s *Node) Format(eachIndent int) string {
	return s.formatTree(0, eachIndent)
}

func (s *Node) formatTree(indent int, eachIndent int) string {
	var builder strings.Builder
	for tmp := s; tmp != nil; tmp = tmp.next {
		builder.WriteString(tmp.formatNode(indent))
		if !tmp.IsTerminal(){
			builder.WriteString(tmp.child.formatTree(indent + eachIndent, eachIndent))
		}
	}
	return builder.String()
}
func (s *Node) formatNode(indent int) string {
	var builder strings.Builder
	builder.WriteString(strings.Repeat(" ", indent))
	switch dat := s.data.(type) {
	case Token:
		switch dat.Kind {
		case Tident:
			builder.WriteString(fmt.Sprintf(" Terminal: %s", dat.Data))
		case Tnumber:
			builder.WriteString(fmt.Sprintf(" Terminal: %d", dat.Data))
		}
	case Rule:
		builder.WriteString(fmt.Sprintf(" Nonterminal: %s", dat.String()))
	default:
		panic("unreachable")
	}
	builder.WriteRune('\n')
	return builder.String()
}
