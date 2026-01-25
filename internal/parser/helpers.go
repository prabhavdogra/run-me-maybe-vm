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
	instructionNumber := int64(0)

	// Validate first token doesn't violate expectations
	nextToken := tokens.PeekToken(1)
	switch tokens[0].Type {
	case token.TypeInt, token.TypeLabel:
		panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error(fmt.Sprintf("program cannot start with a %s reference", tokens[0].Type)))
	case token.TypePush:
		if len(tokens) < 2 || util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeFloat, token.TypeChar, token.TypeString, token.TypeNull) {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error(fmt.Sprintf("expected integer, float, char, or string value after '%s' instruction, but found %s '%s'", tokens[0].Type, nextToken.Type, nextToken.Text)))
		}
		current = current.AddNextNode(tokens[1])
		instructionNumber++
		startIndex++
	case token.TypePushPtr:
		if len(tokens) < 2 || util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeNull) {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error(fmt.Sprintf("expected integer or NULL after 'push_ptr' instruction, but found %s '%s'", nextToken.Type, nextToken.Text)))
		}
		current = current.AddNextNode(tokens[1])
		instructionNumber++
		startIndex++
	case token.TypePushStr:
		if len(tokens) < 2 || nextToken.Type != token.TypeString {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error(fmt.Sprintf("expected string value after 'push_str' instruction, but found %s '%s'", nextToken.Type, nextToken.Text)))
		}
		current = current.AddNextNode(tokens[1])
		// push_str does not increment instructionNumber as it's a data directive
		startIndex++
	case token.TypeGetStr:
		if len(tokens) < 2 || util.NotOneOf(nextToken.Type, token.TypeInt) {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error(fmt.Sprintf("expected integer value (index) after 'get_str' instruction, but found %s '%s'", nextToken.Type, nextToken.Text)))
		}
		current = current.AddNextNode(tokens[1])
		instructionNumber++
		startIndex++
	case token.TypeInDup, token.TypeInSwap, token.TypeInDupStr,
		token.TypeInSwapStr, token.TypeCastIntToFloat, token.TypeCastFloatToInt:
		if len(tokens) < 2 || util.NotOneOf(nextToken.Type, token.TypeInt) {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error(fmt.Sprintf("expected integer value after '%s' instruction, but found %s '%s'", tokens[0].Type, nextToken.Type, nextToken.Text)))
		}
		current = current.AddNextNode(tokens[1])
		instructionNumber++
		startIndex++
	case token.TypeNative:
		if len(tokens) < 2 || nextToken.Type != token.TypeInt {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error("expected integer value (function ID) after 'native' instruction"))
		}
		current = current.AddNextNode(tokens[1])
		instructionNumber++
		startIndex++
	case token.TypeJmp, token.TypeZjmp, token.TypeNzjmp:
		if len(tokens) < 2 {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error("expected label or integer after jump instruction at the start of the program"))
		}
		if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeLabel) {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error("expected label or integer after jump instruction at the start of the program"))
		}
		current = current.AddNextNode(tokens[1])
		instructionNumber++
		startIndex++
	case token.TypeCall:
		if len(tokens) < 2 {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error("expected label or integer after call instruction"))
		}
		if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeLabel) {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error("expected label or integer after call instruction"))
		}
		current = current.AddNextNode(tokens[1])
		instructionNumber++
		startIndex++
	case token.TypeEntrypoint:
		if len(tokens) < 2 {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error("expected label or integer after entrypoint instruction"))
		}
		if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeLabel) {
			panic(token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}.Error("expected label or integer after entrypoint instruction"))
		}
		current = current.AddNextNode(tokens[1])
		startIndex++
	case token.TypeRet:
		instructionNumber++
		startIndex++
	case token.TypeLabelDefinition:
		handleLabelDefination(tokens[0], labelMap, instructionNumber)
		ctx := token.TokenContext{Line: tokens[0].Line, Character: tokens[0].Character, FileName: tokens[0].FileName}
		root = &ParserList{
			Value: token.GetNoOpToken(ctx),
			Next:  nil,
		}
		current = root
		instructionNumber++
	}

	for i := startIndex; i < len(tokens); i++ {
		curToken := tokens[i]
		nextToken := tokens.PeekToken(i + 1)
		switch curToken.Type {
		case token.TypePush:
			if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeFloat, token.TypeChar, token.TypeString, token.TypeNull) {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error(fmt.Sprintf("expected integer, float, char, or string value after '%s' instruction, but found %s '%s'", curToken.Type, nextToken.Type, nextToken.Text)))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			if nextToken.Type == token.TypeString {
				instructionNumber += int64(len(nextToken.Text))
			} else {
				instructionNumber++
			}
			i++
		case token.TypePushPtr:
			if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeNull) {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error(fmt.Sprintf("expected integer or NULL after 'push_ptr' instruction, but found %s '%s'", nextToken.Type, nextToken.Text)))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			instructionNumber++
			i++
		case token.TypePushStr:
			if nextToken.Type != token.TypeString {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error(fmt.Sprintf("expected string value after 'push_str' instruction, but found %s '%s'", nextToken.Type, nextToken.Text)))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			i++
		case token.TypeGetStr:
			if util.NotOneOf(nextToken.Type, token.TypeInt) {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error(fmt.Sprintf("expected integer value (index) after 'get_str' instruction, but found %s '%s'", nextToken.Type, nextToken.Text)))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			instructionNumber++
			i++
		case token.TypeInDup, token.TypeInSwap, token.TypeInDupStr, token.TypeInSwapStr:
			if util.NotOneOf(nextToken.Type, token.TypeInt) {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error(fmt.Sprintf("expected integer value after '%s' instruction, but found %s '%s'", curToken.Type, nextToken.Type, nextToken.Text)))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			instructionNumber++
			i++
		case token.TypeNative:
			if util.NotOneOf(nextToken.Type, token.TypeInt) {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error("expected integer value (function ID) after 'native' instruction"))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			instructionNumber++
			i++
		case token.TypeJmp, token.TypeZjmp, token.TypeNzjmp:
			if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeLabel) {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error(fmt.Sprintf("expected label after '%s' instruction, but found %s '%s'", curToken.Type, nextToken.Type, nextToken.Text)))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			instructionNumber++
			i++
		case token.TypeLabelDefinition:
			ctx := token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}
			current = current.AddNextNode(token.GetNoOpToken(ctx))
			instructionNumber++
			handleLabelDefination(curToken, labelMap, instructionNumber)
		case token.TypeCall:
			if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeLabel) {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error(fmt.Sprintf("expected label or integer after '%s' instruction, but found %s '%s'", curToken.Type, nextToken.Type, nextToken.Text)))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			instructionNumber++
			i++
		case token.TypeEntrypoint:
			if util.NotOneOf(nextToken.Type, token.TypeInt, token.TypeLabel) {
				panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error(fmt.Sprintf("expected label or integer after '%s' instruction, but found %s '%s'", curToken.Type, nextToken.Type, nextToken.Text)))
			}
			current = current.AddNextNode(curToken)
			current = current.AddNextNode(nextToken)
			// No instruction increment
			i++
		case token.TypeRet:
			current = current.AddNextNode(curToken)
			instructionNumber++
		case token.TypeNoOp, token.TypePop, token.TypeDup, token.TypeSwap,
			token.TypeAdd, token.TypeSub, token.TypeMul, token.TypeDiv,
			token.TypeMod, token.TypeCmpe, token.TypeCmpne, token.TypeCmpg,
			token.TypeCmpl, token.TypeCmpge, token.TypeCmple, token.TypePrint,
			token.TypeInt, token.TypeHalt, token.TypeLabel, token.TypeIntToStr,
			token.TypeNull, token.TypePopStr, token.TypeDupStr, token.TypeSwapStr,
			token.TypeCastIntToFloat, token.TypeCastFloatToInt:
			current = current.AddNextNode(curToken)
			instructionNumber++
		default:
			panic(token.TokenContext{Line: curToken.Line, Character: curToken.Character, FileName: curToken.FileName}.Error("unknown token type encountered during parsing"))
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
	length := 0
	for current := pl; current != nil; current = current.Next {
		current.Value.Print()
		length++
		fmt.Println()
	}
	fmt.Println("Parser Length: ", length)
}

func handleLabelDefination(t token.Token, labelMap map[string]int64, instructionNum int64) {
	if _, exists := labelMap[t.Text]; exists {
		panic(token.TokenContext{Line: t.Line, Character: t.Character, FileName: t.FileName}.Error(fmt.Sprintf("duplicate label definition found for label '%s'", t.Text)))
	}
	labelMap[t.Text] = instructionNum
}

func assertAndReplaceLabels(parserList *ParserList, labelMap map[string]int64) {
	cur := parserList
	for cur != nil {
		if cur.Value.Type == token.TypeLabel {
			label := cur.Value.Text
			lineNum, exists := labelMap[label]
			if !exists {
				panic(token.TokenContext{Line: cur.Value.Line, Character: cur.Value.Character, FileName: cur.Value.FileName}.Error(fmt.Sprintf("undefined label reference found for label '%s'", label)))
			}
			// Replace label token with integer token representing the instruction number
			cur.Value.Type = token.TypeInt
			cur.Value.Text = fmt.Sprintf("%d", lineNum)
		}
		cur = cur.Next
	}
}
