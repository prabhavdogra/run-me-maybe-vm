package parser

import (
	"vm/internal/lexer"
	"vm/internal/token"
)

func Init(l *lexer.Lexer) *ParserList {
	if l.Tokens == nil || len(l.Tokens) == 0 {
		return nil
	}

	parserList := generateList(l.Tokens)
	return parserList
}

func generateList(tokens token.Tokens) *ParserList {
	if len(tokens) == 0 {
		return nil
	}

	// Start parser list with the first token
	root := &ParserList{
		Value: tokens[0],
		Next:  nil,
	}

	// Validate first token doesn't violate expectations
	switch tokens[0].Type {
	case token.TypeInt:
		panic("ERROR: program cannot start with an integer token")
	case token.TypePush, token.TypeInDup, token.TypeInSwap:
		if len(tokens) < 2 || tokens[1].Type != token.TypeInt {
			panic("ERROR: expected integer value after push/indup/inswap instruction at the start of the program")
		}
	}

	// Iterate through tokens by absolute index and append to the parser list.
	// When we encounter instructions that require an integer operand, also
	// append that integer token and advance the index to consume it.
	for i := 1; i < len(tokens); i++ {
		t := tokens[i]
		switch t.Type {
		case token.TypePush:
			if tokens.PeekToken(i+1).Type != token.TypeInt {
				panic("ERROR: expected integer value after push instruction")
			}
			root.AddNextNode(t)
			root.AddNextNode(tokens[i+1])
			i++ // consume the integer token
		case token.TypeInDup:
			if tokens.PeekToken(i+1).Type != token.TypeInt {
				panic("ERROR: expected integer value after indup instruction")
			}
			root.AddNextNode(t)
			root.AddNextNode(tokens[i+1])
			i++
		case token.TypeInSwap:
			if tokens.PeekToken(i+1).Type != token.TypeInt {
				panic("ERROR: expected integer value after inswap instruction")
			}
			root.AddNextNode(t)
			root.AddNextNode(tokens[i+1])
			i++
		case token.TypeNoOp, token.TypePop, token.TypeDup, token.TypeSwap,
			token.TypeAdd, token.TypeSub, token.TypeMul, token.TypeDiv,
			token.TypeMod, token.TypeCmpe, token.TypeCmpne, token.TypeCmpg,
			token.TypeCmpl, token.TypeCmpge, token.TypeCmple, token.TypeJmp,
			token.TypeZjmp, token.TypeNzjmp, token.TypePrint, token.TypeInt, token.TypeHalt:
			root.AddNextNode(t)
		default:
			panic("unknown token type encountered during parsing")
		}
	}

	return root
}

func (pl *ParserList) AddNextNode(token token.Token) {
	current := pl
	for current.Next != nil {
		current = current.Next
	}
	current.Next = &ParserList{
		Value: token,
		Next:  nil,
	}
}

func (pl *ParserList) Print() {
	for current := pl; current != nil; current = current.Next {
		current.Value.Print()
	}
}
