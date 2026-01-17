package parser

import "vm/internal/token"

// Parser represents linked list of tokens
type ParserList struct {
	Value token.Token
	Next  *ParserList
}
