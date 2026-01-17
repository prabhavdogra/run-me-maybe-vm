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
	root := &ParserList{
		Value: tokens[0],
		Next:  nil,
	}
	for id, parsedToken := range tokens {
		switch parsedToken.Type {
		case token.TypeNoOp:
			root.AddNextNode(parsedToken)
		case token.TypePush:
			root.AddNextNode(parsedToken)
			if tokens.PeekToken(id+1).Type != token.TypeInt {
				panic("ERROR: expected integer value after push instruction")
			}
		case token.TypePop:
			root.AddNextNode(parsedToken)
		case token.TypeDup:
			root.AddNextNode(parsedToken)
		case token.TypeIndup:
			root.AddNextNode(parsedToken)
			if tokens.PeekToken(id+1).Type != token.TypeInt {
				panic("ERROR: expected integer value after indup instruction")
			}
		case token.TypeSwap:
			root.AddNextNode(parsedToken)
		case token.TypeInswap:
			root.AddNextNode(parsedToken)
			if tokens.PeekToken(id+1).Type != token.TypeInt {
				panic("ERROR: expected integer value after inswap instruction")
			}
		case token.TypeAdd:
			root.AddNextNode(parsedToken)
		case token.TypeSub:
			root.AddNextNode(parsedToken)
		case token.TypeMul:
			root.AddNextNode(parsedToken)
		case token.TypeDiv:
			root.AddNextNode(parsedToken)
		case token.TypeMod:
			root.AddNextNode(parsedToken)
		case token.TypeCmpe:
			root.AddNextNode(parsedToken)
		case token.TypeCmpne:
			root.AddNextNode(parsedToken)
		case token.TypeCmpg:
			root.AddNextNode(parsedToken)
		case token.TypeCmpl:
			root.AddNextNode(parsedToken)
		case token.TypeCmpge:
			root.AddNextNode(parsedToken)
		case token.TypeCmple:
			root.AddNextNode(parsedToken)
		case token.TypeJmp:
			root.AddNextNode(parsedToken)
		case token.TypeZjmp:
			root.AddNextNode(parsedToken)
		case token.TypeNzjmp:
			root.AddNextNode(parsedToken)
		case token.TypePrint:
			root.AddNextNode(parsedToken)
		case token.TypeInt:
			root.AddNextNode(parsedToken)
		case token.TypeHalt:
			root.AddNextNode(parsedToken)
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
