package lexer

import (
	"fmt"
	"os"
	"unicode"
)

func printToken(token Token) {
	switch token.Type {
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
	}

	fmt.Printf(
		"text: %s, line: %d, character: %d\n",
		token.Text,
		token.Line,
		token.Character,
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
		return TypeNoOp
	}
}

func openFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s:  %w", filePath, err)
	}
	return data, nil
}

func generateKeyword(input string, line int, currentIndex int, char int) (Token, int) {
	keyword := ""
	for len(input) > currentIndex && unicode.IsLetter(rune(input[currentIndex])) {
		keyword += string(input[currentIndex])
		currentIndex++
	}
	tokenType := checkBuiltinKeywords(keyword)
	return InitToken(tokenType, keyword, line, char), currentIndex
}

func Lex() {
	byteArray, err := openFile("test.wm")
	if err != nil {
		fmt.Println(err)
		return
	}

	input := string(byteArray)
	currentIndex := 0
	line := 1
	character := 1

	for currentIndex < len(input) {

		if input[currentIndex] == '\n' {
			line++
			character = 0
		}

		if unicode.IsLetter(rune(input[currentIndex])) {
			var token Token
			token, currentIndex = generateKeyword(input, line, currentIndex, character)
			printToken(token)
		} else if unicode.IsDigit(rune(input[currentIndex])) {
			fmt.Println("NUMERIC")
			currentIndex++
		} else { // whitespace token
			currentIndex++
		}
		character++
	}

}
