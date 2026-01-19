package parser

import (
	"fmt"
	"vm/internal/lexer"
	"vm/internal/token"
	"vm/util"
)

func Init(l *lexer.Lexer) *ParserList {
	if l.Tokens == nil || len(l.Tokens) == 0 {
		return nil
	}
	labelMap := make(map[string]int64)
	parserList := generateList(l.Tokens, labelMap)
	return parserList
}

func generateList(tokens token.Tokens, labelMap map[string]int64) *ParserList {
	if len(tokens) == 0 {
		return nil
	}

	// Start parser list with the first token
	root := &ParserList{
		Value: tokens[0],
		Next:  nil,
	}

	current := root
	startIndex := 1

	// Validate first token doesn't violate expectations
	nextToken := tokens.PeekToken(1)
	switch tokens[0].Type {
	case token.TypeInt, token.TypeLabel:
		instructionType := tokens[0].Type
		panic(fmt.Sprintf("ERROR: program cannot start with a %s reference", instructionType))
	case token.TypePush, token.TypeInDup, token.TypeInSwap:
		if len(tokens) < 2 || util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeFloat) {
			instructionType := tokens[0].Type
			panic(fmt.Sprintf("ERROR: expected integer token after %s instruction", instructionType))
		}
		current = current.AddNextNode(tokens[1])
		startIndex++
	case token.TypeJmp, token.TypeZjmp, token.TypeNzjmp:
		if len(tokens) < 2 {
			panic("ERROR: expected label or integer after jump instruction at the start of the program")
		}
		if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeLabel) {
			panic("ERROR: expected label or integer after jump instruction at the start of the program")
		}
		current = current.AddNextNode(tokens[1])
		startIndex++
	case token.TypeLabelDefinition:
		handleLabelDefination(tokens[0], labelMap)
		root = &ParserList{
			Value: token.GetNoOpToken(tokens[0].Line, tokens[0].Character),
			Next:  nil,
		}
	}

	// Iterate through tokens by absolute index and append to the parser list.
	// When we encounter instructions that require an integer operand, also
	// append that integer token and advance the index to consume it.
	for i := startIndex; i < len(tokens); i++ {
		curToken := tokens[i]
		nextToken := tokens.PeekToken(i + 1)
		switch curToken.Type {
		case token.TypePush, token.TypeInDup, token.TypeInSwap:
			if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeFloat) {
				instructionType := curToken.Type.String()
				panic(fmt.Sprintf("ERROR: expected integer token after %s instruction", instructionType))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			i++
		case token.TypeJmp, token.TypeZjmp, token.TypeNzjmp:
			if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeLabel) {
				instructionType := curToken.Type.String()
				panic(fmt.Sprintf("ERROR: expected label token after %s instruction", instructionType))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			i++
		case token.TypeLabelDefinition:
			handleLabelDefination(curToken, labelMap)
			current = current.AddNextNode(token.GetNoOpToken(curToken.Line, curToken.Character))
		case token.TypeNoOp, token.TypePop, token.TypeDup, token.TypeSwap,
			token.TypeAdd, token.TypeSub, token.TypeMul, token.TypeDiv,
			token.TypeMod, token.TypeCmpe, token.TypeCmpne, token.TypeCmpg,
			token.TypeCmpl, token.TypeCmpge, token.TypeCmple, token.TypePrint,
			token.TypeInt, token.TypeHalt, token.TypeLabel:
			current = current.AddNextNode(curToken)
		default:
			panic("unknown token type encountered during parsing")
		}
	}

	assertAndReplaceLabels(root, labelMap)

	return root
}

func (pl *ParserList) AddNextNode(token token.Token) *ParserList {
	pl.Next = &ParserList{
		Value: token,
		Next:  nil,
	}
	pl = pl.Next
	return pl
}

func (pl *ParserList) Print() {
	for current := pl; current != nil; current = current.Next {
		current.Value.Print()
	}
}

func handleLabelDefination(t token.Token, labelMap map[string]int64) {
	if _, exists := labelMap[t.Text]; exists {
		panic(fmt.Sprintf("ERROR: duplicate label definition found for label '%s'", t.Text))
	}
	labelMap[t.Text] = t.Line
}

func assertAndReplaceLabels(parserList *ParserList, labelMap map[string]int64) {
	cur := parserList
	for cur != nil {
		if cur.Value.Type == token.TypeLabel {
			label := cur.Value.Text
			lineNum, exists := labelMap[label]
			if !exists {
				panic(fmt.Sprintf("ERROR: undefined label reference found for label '%s'", label))
			}
			// Replace label token with integer token representing the line number
			cur.Value.Type = token.TypeInt
			cur.Value.Text = fmt.Sprintf("%d", lineNum)
		}
		cur = cur.Next
	}
}
