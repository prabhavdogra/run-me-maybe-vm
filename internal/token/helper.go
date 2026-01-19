package token

import (
	"fmt"
	"unicode"
)

type Tokens []Token

func (tt TokenType) String() string {
	switch tt {
	case TypeNoOp:
		return "TYPE NOP"
	case TypePush:
		return "TYPE PUSH"
	case TypePop:
		return "TYPE POP"
	case TypeDup:
		return "TYPE DUP"
	case TypeInDup:
		return "TYPE INDUP"
	case TypeSwap:
		return "TYPE SWAP"
	case TypeInSwap:
		return "TYPE INSWAP"
	case TypeAdd:
		return "TYPE ADD"
	case TypeSub:
		return "TYPE SUB"
	case TypeMul:
		return "TYPE MUL"
	case TypeDiv:
		return "TYPE DIV"
	case TypeMod:
		return "TYPE MOD"
	case TypeCmpe:
		return "TYPE CMPE"
	case TypeCmpne:
		return "TYPE CMPNE"
	case TypeCmpg:
		return "TYPE CMPG"
	case TypeCmpl:
		return "TYPE CMPL"
	case TypeCmpge:
		return "TYPE CMPGE"
	case TypeCmple:
		return "TYPE CMPLE"
	case TypeJmp:
		return "TYPE JMP"
	case TypeZjmp:
		return "TYPE ZJMP"
	case TypeNzjmp:
		return "TYPE NZJMP"
	case TypePrint:
		return "TYPE PRINT"
	case TypeInt:
		return "TYPE INT"
	case TypeFloat:
		return "TYPE FLOAT"
	case TypeChar:
		return "TYPE CHAR"
	case TypeLabelDefinition:
		return "TYPE LABEL DEFINITION"
	case TypeLabel:
		return "TYPE LABEL"
	case TypeHalt:
		return "TYPE HALT"
	default:
		return "TYPE INVALID"
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

func InitToken(tokenType TokenType, text string, line int64, char int) Token {
	return Token{
		Type:      tokenType,
		Text:      text,
		Line:      line,
		Character: char,
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

func GenerateKeyword(input string, line int64, currentIndex int, char int) (Token, int) {
	keyword := ""
	for len(input) > currentIndex &&
		(unicode.IsLetter(rune(input[currentIndex])) ||
			rune(input[currentIndex]) == ':' ||
			rune(input[currentIndex]) == '_') {
		keyword += string(rune(input[currentIndex]))
		currentIndex++
	}
	tokenType := checkBuiltinKeywords(keyword)
	if tokenType == TypeLabelDefinition {
		keyword = keyword[:len(keyword)-1]
	}
	return InitToken(tokenType, keyword, line, char), currentIndex
}

func GenerateNumber(input string, line int64, currentIndex int, char int) (Token, int) {
	number := ""
	for len(input) > currentIndex && unicode.IsDigit(rune(input[currentIndex])) {
		number += string(rune(input[currentIndex]))
		currentIndex++
	}
	// Integer case
	if len(input) <= currentIndex || input[currentIndex] != '.' {
		return InitToken(TypeInt, number, line, char), currentIndex
	}
	number = number + string(input[currentIndex])
	currentIndex++
	// Float case
	for len(input) > currentIndex && unicode.IsDigit(rune(input[currentIndex])) {
		number += string(input[currentIndex])
		currentIndex++
	}

	return InitToken(TypeFloat, number, line, char), currentIndex
}

func GenerateChar(input string, line int64, currentIndex int, char int) (Token, int) {
	currentIndex++ // skip opening '

	if currentIndex >= len(input) {
		panic(fmt.Sprintf("ERROR: unterminated character literal at line %d", line))
	}

	charValue := input[currentIndex]
	currentIndex++ // skip the character

	if currentIndex >= len(input) || input[currentIndex] != '\'' {
		panic(fmt.Sprintf("ERROR: unterminated character literal at line %d", line))
	}

	currentIndex++ // skip closing '

	return InitToken(TypeChar, string(charValue), line, char), currentIndex
}

func (t Tokens) PeekToken(index int) Token {
	if index < 0 || index >= len(t) {
		return Token{
			Type: TypeInvalid,
		}
	}
	return t[index]
}

func GetNoOpToken(line int64, char int) Token {
	return InitToken(TypeNoOp, "noop", line, char)
}
