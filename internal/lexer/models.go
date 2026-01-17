package lexer

import "vm/internal/token"

type Lexer struct {
	Tokens   token.Tokens
	fileName string
}
