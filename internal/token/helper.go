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
	case TypePushPtr:
		return "push_ptr"
	case TypePushStr:
		return "push_str"
	case TypeGetStr:
		return "get_str"
	case TypePopStr:
		return "pop_str"
	case TypeDupStr:
		return "dup_str"
	case TypeInDupStr:
		return "indup_str"
	case TypeSwapStr:
		return "swap_str"
	case TypeInSwapStr:
		return "inswap_str"
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
	case TypeNative:
		return "native"
	case TypeInt:
		return "int"
	case TypeFloat:
		return "float"
	case TypeChar:
		return "char"
	case TypeString:
		return "string"
	case TypeLabelDefinition:
		return "label definition"
	case TypeLabel:
		return "label"
	case TypeHalt:
		return "halt"
	case TypeIntToStr:
		return "int_to_str"
	case TypeNull:
		return "NULL"
	case TypeCall:
		return "call"
	case TypeRet:
		return "ret"
	case TypeEntrypoint:
		return "entrypoint"
	case TypeCastIntToFloat:
		return "itof"
	case TypeCastFloatToInt:
		return "ftoi"
	case TypeRef:
		return "ref"
	case TypeDeref:
		return "deref"
	case TypeMovStr:
		return "mov_str"
	case TypeIndex:
		return "index"
	case TypeRegister:
		return "register"
	case TypeMov:
		return "mov"
	case TypeTop:
		return "top"
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
	case "push_str":
		return TypePushStr
	case "get_str":
		return TypeGetStr
	case "pop_str":
		return TypePopStr
	case "dup_str":
		return TypeDupStr
	case "indup_str":
		return TypeInDupStr
	case "swap_str":
		return TypeSwapStr
	case "inswap_str":
		return TypeInSwapStr
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
	case "push_ptr":
		return TypePushPtr
	case "native":
		return TypeNative
	case "halt":
		return TypeHalt
	case "int_to_str":
		return TypeIntToStr
	case "NULL":
		return TypeNull
	case "call":
		return TypeCall
	case "ret":
		return TypeRet
	case "entrypoint":
		return TypeEntrypoint
	case "itof":
		return TypeCastIntToFloat
	case "ftoi":
		return TypeCastFloatToInt
	case "ref":
		return TypeRef
	case "deref":
		return TypeDeref
	case "mov_str":
		return TypeMovStr
	case "index":
		return TypeIndex
	case "mov":
		return TypeMov
	case "top":
		return TypeTop
	default:
		return checkLabelType(name)
	}
}

func GetWord(input string, currentIndex int) (string, int) {
	keyword := ""
	for len(input) > currentIndex &&
		(unicode.IsLetter(rune(input[currentIndex])) ||
			unicode.IsDigit(rune(input[currentIndex])) ||
			rune(input[currentIndex]) == ':' ||
			rune(input[currentIndex]) == '_') {
		keyword += string(rune(input[currentIndex]))
		currentIndex++
	}
	return keyword, currentIndex
}

func checkRegisterType(name string) TokenType {
	if len(name) > 1 && name[0] == 'r' {
		for _, r := range name[1:] {
			if !unicode.IsDigit(r) {
				return TypeInvalid
			}
		}
		return TypeRegister
	}
	return TypeInvalid
}

func GenerateKeyword(input string, currentIndex int, ctx TokenContext, macros map[string]string) (Token, string, int) {
	keyword, updatedIndex := GetWord(input, currentIndex)
	if val, ok := macros[keyword]; ok {
		return Token{}, val, updatedIndex
	}
	tokenType := checkRegisterType(keyword)
	if tokenType == TypeInvalid {
		tokenType = checkBuiltinKeywords(keyword)
	}

	if tokenType == TypeLabelDefinition {
		keyword = keyword[:len(keyword)-1]
	}
	return InitToken(tokenType, keyword, ctx), "", updatedIndex
}

func GenerateNumber(input string, currentIndex int, ctx TokenContext) (Token, int) {
	number := ""
	if len(input) > currentIndex && input[currentIndex] == '-' {
		number += "-"
		currentIndex++
	}
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
	if charValue == '\\' {
		currentIndex++
		if currentIndex >= len(input) {
			panic(fmt.Sprintf("ERROR: unterminated character literal at line %d", ctx.Line))
		}
		escapeChar := input[currentIndex]
		switch escapeChar {
		case 'n':
			charValue = '\n'
		case 't':
			charValue = '\t'
		case 'r':
			charValue = '\r'
		case '\\':
			charValue = '\\'
		case '\'':
			charValue = '\''
		case '0':
			charValue = 0
		default:
			panic(fmt.Sprintf("ERROR: unknown escape character '\\%c' at line %d", escapeChar, ctx.Line))
		}
	}
	currentIndex++ // skip the character (or the escape code)

	if currentIndex >= len(input) || input[currentIndex] != '\'' {
		panic(fmt.Sprintf("ERROR: unterminated character literal at line %d", ctx.Line))
	}

	currentIndex++ // skip closing '

	return InitToken(TypeChar, string(charValue), ctx), currentIndex
}

func GenerateString(input string, currentIndex int, ctx TokenContext) (Token, int) {
	currentIndex++ // skip opening "
	if currentIndex >= len(input) {
		panic(ctx.Error("unterminated string literal"))
	}

	var strValue string
	for currentIndex < len(input) && input[currentIndex] != '"' {
		if input[currentIndex] == '\\' { // escape character
			currentIndex++ // skip backslash
			if currentIndex >= len(input) {
				panic(ctx.Error("unterminated string literal"))
			}
			switch input[currentIndex] {
			case 'n':
				strValue += "\n"
			case 't':
				strValue += "\t"
			case '"':
				strValue += "\""
			case '\\':
				strValue += "\\"
			case '0':
				strValue += "\000"
			default:
				panic(ctx.Error(fmt.Sprintf("unknown escape character: \\%c", input[currentIndex])))
			}
		} else {
			strValue += string(input[currentIndex])
		}
		currentIndex++
	}

	if currentIndex >= len(input) {
		panic(ctx.Error("unterminated string literal"))
	}

	currentIndex++ // skip closing "
	return InitToken(TypeString, strValue, ctx), currentIndex
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
