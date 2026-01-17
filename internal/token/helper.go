package token

import (
	"fmt"
	"unicode"
)

type Tokens []Token

func (t *Token) Print() {
	switch t.Type {
	case TypeNoOp:
		fmt.Println("TYPE NOP")
	case TypePush:
		fmt.Println("TYPE PUSH")
	case TypePop:
		fmt.Println("TYPE POP")
	case TypeDup:
		fmt.Println("TYPE DUP")
	case TypeIndup:
		fmt.Println("TYPE INDUP")
	case TypeSwap:
		fmt.Println("TYPE SWAP")
	case TypeInswap:
		fmt.Println("TYPE INSWAP")
	case TypeAdd:
		fmt.Println("TYPE ADD")
	case TypeSub:
		fmt.Println("TYPE SUB")
	case TypeMul:
		fmt.Println("TYPE MUL")
	case TypeDiv:
		fmt.Println("TYPE DIV")
	case TypeMod:
		fmt.Println("TYPE MOD")
	case TypeCmpe:
		fmt.Println("TYPE CMPE")
	case TypeCmpne:
		fmt.Println("TYPE CMPNE")
	case TypeCmpg:
		fmt.Println("TYPE CMPG")
	case TypeCmpl:
		fmt.Println("TYPE CMPL")
	case TypeCmpge:
		fmt.Println("TYPE CMPGE")
	case TypeCmple:
		fmt.Println("TYPE CMPLE")
	case TypeJmp:
		fmt.Println("TYPE JMP")
	case TypeZjmp:
		fmt.Println("TYPE ZJMP")
	case TypeNzjmp:
		fmt.Println("TYPE NZJMP")
	case TypePrint:
		fmt.Println("TYPE PRINT")
	case TypeHalt:
		fmt.Println("TYPE HALT")
	case TypeInt:
		fmt.Println("TYPE INT")
	default:
		fmt.Println("TYPE INVALID")
	}

	fmt.Printf(
		"text: %s, line: %d, character: %d\n",
		t.Text,
		t.Line,
		t.Character,
	)
}

func InitToken(tokenType TokenType, text string, line int, char int) Token {
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
		return TypeIndup
	case "swap":
		return TypeSwap
	case "inswap":
		return TypeInswap
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
		return TypeInvalid
	}
}

func GenerateKeyword(input string, line int, currentIndex int, char int) (Token, int) {
	keyword := ""
	for len(input) > currentIndex && unicode.IsLetter(rune(input[currentIndex])) {
		keyword += string(input[currentIndex])
		currentIndex++
	}
	tokenType := checkBuiltinKeywords(keyword)
	return InitToken(tokenType, keyword, line, char), currentIndex
}

func GenerateInt(input string, line int, currentIndex int, char int) (Token, int) {
	number := ""
	for len(input) > currentIndex && unicode.IsDigit(rune(input[currentIndex])) {
		number += string(input[currentIndex])
		currentIndex++
	}
	return InitToken(TypeInt, number, line, char), currentIndex
}

func (t Tokens) PeekToken(index int) *Token {
	if index < 0 || index >= len(t) {
		return nil
	}
	return &t[index]
}
