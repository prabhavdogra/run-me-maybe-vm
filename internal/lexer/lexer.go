package lexer

import (
	"fmt"
	"os"
	"unicode"
	"vm/internal/token"
)

func (l *Lexer) openFile() (string, error) {
	data, err := os.ReadFile(l.fileName)
	if err != nil {
		return "", fmt.Errorf("could not open file %s:  %w", l.fileName, err)
	}
	return string(data), nil
}

func Init(filename string) *Lexer {
	return &Lexer{
		Tokens:   []token.Token{},
		fileName: filename,
	}
}

func (l *Lexer) addToken(token token.Token) *Lexer {
	l.Tokens = append(l.Tokens, token)
	return l
}

func (l *Lexer) Lex() *Lexer {
	input, err := l.openFile()
	if err != nil {
		fmt.Println(err)
		return l
	}

	currentIndex := 0
	line := int64(1)
	character := 1

	for currentIndex < len(input) {
		var lexedToken token.Token
		if input[currentIndex] == '\n' {
			if (currentIndex == 0 && len(l.Tokens) == 0) ||
				(currentIndex != 0 && input[currentIndex-1] == '\n') {
				l.addToken(token.GetNoOpToken(line, character))
			}
			line++
			character = 0
		}

		if unicode.IsLetter(rune(input[currentIndex])) { // keyword token
			lexedToken, currentIndex = token.GenerateKeyword(input, line, currentIndex, character)
			l.addToken(lexedToken)
		} else if unicode.IsDigit(rune(input[currentIndex])) { // numeric token
			lexedToken, currentIndex = token.GenerateNumber(input, line, currentIndex, character)
			l.addToken(lexedToken)
		} else if input[currentIndex] == '\'' { // character literal
			lexedToken, currentIndex = token.GenerateChar(input, line, currentIndex, character)
			l.addToken(lexedToken)
		} else { // whitespace token
			currentIndex++
		}
		character++
	}
	return l
}

func (l *Lexer) Print() {
	for _, token := range l.Tokens {
		token.Print()
	}
}
