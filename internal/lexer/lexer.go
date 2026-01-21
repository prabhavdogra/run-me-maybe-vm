package lexer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"vm/internal/token"
)

type Lexer struct {
	Tokens   []token.Token
	FileName string
	Macros   map[string]string
}

func Init(filename string) *Lexer {
	return &Lexer{
		Tokens:   []token.Token{},
		FileName: filename,
		Macros:   make(map[string]string),
	}
}

func (l *Lexer) addToken(token token.Token) *Lexer {
	l.Tokens = append(l.Tokens, token)
	return l
}

func (l *Lexer) Lex() *Lexer {
	l.processFile(l.FileName)
	return l
}

func (l *Lexer) processFile(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(fmt.Errorf("could not open file %s:  %w", fileName, err))
	}
	l.lexContent(string(data), fileName)
}

func (l *Lexer) processImport(ctx *token.TokenContext, input string, currentIndex int) int {
	for currentIndex < len(input) && unicode.IsSpace(rune(input[currentIndex])) {
		if input[currentIndex] == '\n' {
			ctx.Line++
			ctx.Character = 0
		}
		currentIndex++
	}
	// expect quote
	if currentIndex >= len(input) || input[currentIndex] != '"' {
		panic(fmt.Sprintf("ERROR (%s:%d): expected filename in quotes after @imp", ctx.FileName, ctx.Line))
	}
	currentIndex++ // skip quote
	importFile := ""
	for currentIndex < len(input) && input[currentIndex] != '"' {
		importFile += string(input[currentIndex])
		currentIndex++
	}
	if currentIndex >= len(input) {
		panic(fmt.Sprintf("ERROR (%s:%d): unterminated string in @imp", ctx.FileName, ctx.Line))
	}
	currentIndex++ // skip closing quote
	// Resolve import path relative to current file's directory
	importPath := importFile
	if !filepath.IsAbs(importFile) {
		currentDir := filepath.Dir(ctx.FileName)
		importPath = filepath.Join(currentDir, importFile)
	}
	l.processFile(importPath)
	return currentIndex
}

func (l *Lexer) processDef(ctx *token.TokenContext, input string, currentIndex int) int {
	// skip whitespace
	for currentIndex < len(input) && unicode.IsSpace(rune(input[currentIndex])) && input[currentIndex] != '\n' {
		currentIndex++
	}
	// get key
	key := ""
	for currentIndex < len(input) && !unicode.IsSpace(rune(input[currentIndex])) {
		key += string(input[currentIndex])
		currentIndex++
	}
	// skip whitespace
	for currentIndex < len(input) && unicode.IsSpace(rune(input[currentIndex])) && input[currentIndex] != '\n' {
		currentIndex++
	}
	// get value (until newline)
	val := ""
	for currentIndex < len(input) && input[currentIndex] != '\n' {
		val += string(input[currentIndex])
		currentIndex++
	}
	if _, exists := l.Macros[key]; exists {
		panic(ctx.Error(fmt.Sprintf("duplicate macro definition found for macro '%s'", key)))
	}
	l.Macros[key] = strings.TrimSpace(val)
	return currentIndex
}

func (l *Lexer) lexContent(input string, fileName string) {
	currentIndex := 0
	line := int64(1)
	character := 1

	for currentIndex < len(input) {
		var lexedToken token.Token

		// Preprocessor directives
		if input[currentIndex] == '@' {
			currentIndex++ // skip '@'
			directive := ""
			directive, currentIndex = token.GetWord(input, currentIndex)
			ctx := token.TokenContext{Line: line, Character: character, FileName: fileName}
			switch directive {
			case "imp": // @imp
				currentIndex = l.processImport(&ctx, input, currentIndex)
			case "def": // @def
				currentIndex = l.processDef(&ctx, input, currentIndex)
			default:
				panic(fmt.Sprintf("ERROR (%s:%d): checking for unknown preprocessor directive @%s", fileName, line, directive))
			}
			continue
		}
		ctx := token.TokenContext{Line: line, Character: character, FileName: fileName}
		if input[currentIndex] == ';' {
			for currentIndex < len(input) && input[currentIndex] != '\n' {
				currentIndex++
			}
			l.addToken(token.GetNoOpToken(ctx))
		} else if input[currentIndex] == '\n' {
			if (currentIndex == 0) || (input[currentIndex-1] == '\n') {
				l.addToken(token.GetNoOpToken(ctx))
			}
			line++
			character = 0
			currentIndex++
		} else if unicode.IsLetter(rune(input[currentIndex])) { // keyword or macro
			var macroVal string
			lexedToken, macroVal, currentIndex = token.GenerateKeyword(input, currentIndex, ctx, l.Macros)
			if macroVal != "" {
				l.lexContent(macroVal, fileName)
			} else {
				l.addToken(lexedToken)
			}
		} else if unicode.IsDigit(rune(input[currentIndex])) { // numeric token
			lexedToken, currentIndex = token.GenerateNumber(input, currentIndex, ctx)
			l.addToken(lexedToken)
		} else if input[currentIndex] == '\'' { // character literal
			lexedToken, currentIndex = token.GenerateChar(input, currentIndex, ctx)
			l.addToken(lexedToken)
		} else if input[currentIndex] == '"' { // string literal
			lexedToken, currentIndex = token.GenerateString(input, currentIndex, ctx)
			l.addToken(lexedToken)
		} else { // whitespace token
			currentIndex++
		}
		character++
	}
}

func (l *Lexer) Print() {
	for _, token := range l.Tokens {
		token.Print()
		fmt.Println()
	}
}
