package token

import (
	"fmt"
	"unicode"
)

type Tokens []Token

func (tt TokenType) String() string {
	switch tt {
	case TypeNoOp:
		return "noop"
	case TypePush:
		return "push"
	case TypePop:
		return "pop"
	case TypeDup:
		return "dup"
	case TypeInDup:
		return "indup"
	case TypeSwap:
		return "swap"
	case TypeInSwap:
		return "inswap"
	case TypeAdd:
		return "add"
	case TypeSub:
		return "sub"
	case TypeMul:
		return "mul"
	case TypeDiv:
		return "div"
	case TypeMod:
		return "mod"
	case TypeCmpe:
		return "cmpe"
	case TypeCmpne:
		return "cmpne"
	case TypeCmpg:
		return "cmpg"
	case TypeCmpl:
		return "cmpl"
	case TypeCmpge:
		return "cmpge"
	case TypeCmple:
		return "cmple"
	case TypeJmp:
		return "jmp"
	case TypeZjmp:
		return "zjmp"
	case TypeNzjmp:
		return "nzjmp"
	case TypePrint:
		return "print"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeChar:
		return "char"
	case TypeLabelDefinition:
		return "label definition"
	case TypeLabel:
		return "label"
	case TypeHalt:
		return "halt"
	default:
		return "invalid"
	}
}

func (t *Token) Print() {
	fmt.Printf(
		"type: %s text: %s, ",
		t.Type.String(),
		t.Text,
	)
}

func checkLabelType(label string) TokenType {
	if len(label) >= 2 && label[len(label)-1] == ':' {
		return TypeLabelDefinition
	}
	return TypeLabel
}

func InitToken(tokenType TokenType, text string, ctx TokenContext) Token {
	return Token{
		Type:      tokenType,
		Text:      text,
		Line:      ctx.Line,
		Character: ctx.Character,
		FileName:  ctx.FileName,
	}
}

func checkBuiltinKeywords(name string) TokenType {
	switch name {
	case "noop":
		return TypeNoOp
	case "push":
		return TypePush
	case "pop":
		return TypePop
	case "dup":
		return TypeDup
	case "indup":
		return TypeInDup
	case "swap":
		return TypeSwap
	case "inswap":
		return TypeInSwap
	case "add":
		return TypeAdd
	case "sub":
		return TypeSub
	case "mul":
		return TypeMul
	case "div":
		return TypeDiv
	case "mod":
		return TypeMod
	case "cmpe":
		return TypeCmpe
	case "cmpne":
		return TypeCmpne
	case "cmpg":
		return TypeCmpg
	case "cmpl":
		return TypeCmpl
	case "cmpge":
		return TypeCmpge
	case "cmple":
		return TypeCmple
	case "jmp":
		return TypeJmp
	case "zjmp":
		return TypeZjmp
	case "nzjmp":
		return TypeNzjmp
	case "print":
		return TypePrint
	case "halt":
		return TypeHalt
	default:
		return checkLabelType(name)
	}
}

func GetWord(input string, currentIndex int) (string, int) {
	keyword := ""
	for len(input) > currentIndex &&
		(unicode.IsLetter(rune(input[currentIndex])) ||
			rune(input[currentIndex]) == ':' ||
			rune(input[currentIndex]) == '_') {
		keyword += string(rune(input[currentIndex]))
		currentIndex++
	}
	return keyword, currentIndex
}

func GenerateKeyword(input string, currentIndex int, ctx TokenContext, macros map[string]string) (Token, string, int) {
	keyword, updatedIndex := GetWord(input, currentIndex)
	if val, ok := macros[keyword]; ok {
		return Token{}, val, updatedIndex
	}
	tokenType := checkBuiltinKeywords(keyword)
	if tokenType == TypeLabelDefinition {
		keyword = keyword[:len(keyword)-1]
	}
	return InitToken(tokenType, keyword, ctx), "", updatedIndex
}

func GenerateNumber(input string, currentIndex int, ctx TokenContext) (Token, int) {
	number := ""
	for len(input) > currentIndex && unicode.IsDigit(rune(input[currentIndex])) {
		number += string(rune(input[currentIndex]))
		currentIndex++
	}
	// Integer case
	if len(input) <= currentIndex || input[currentIndex] != '.' {
		return InitToken(TypeInt, number, ctx), currentIndex
	}
	number = number + string(input[currentIndex])
	currentIndex++
	// Float case
	for len(input) > currentIndex && unicode.IsDigit(rune(input[currentIndex])) {
		number += string(input[currentIndex])
		currentIndex++
	}

	return InitToken(TypeFloat, number, ctx), currentIndex
}

func GenerateChar(input string, currentIndex int, ctx TokenContext) (Token, int) {
	currentIndex++ // skip opening '

	if currentIndex >= len(input) {
		panic(fmt.Sprintf("ERROR: unterminated character literal at line %d", ctx.Line))
	}

	charValue := input[currentIndex]
	currentIndex++ // skip the character

	if currentIndex >= len(input) || input[currentIndex] != '\'' {
		panic(fmt.Sprintf("ERROR: unterminated character literal at line %d", ctx.Line))
	}

	currentIndex++ // skip closing '

	return InitToken(TypeChar, string(charValue), ctx), currentIndex
}

func (t Tokens) PeekToken(index int) Token {
	if index < 0 || index >= len(t) {
		return Token{
			Type: TypeInvalid,
		}
	}
	return t[index]
}

func GetNoOpToken(ctx TokenContext) Token {
	return InitToken(TypeNoOp, "noop", ctx)
}
