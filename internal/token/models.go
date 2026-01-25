package token

import "fmt"

// TokenContext contains metadata for token creation
type TokenContext struct {
	Line      int64
	Character int
	FileName  string
}

// Error formats an error message with file and line context
func (ctx TokenContext) Error(message string) string {
	return fmt.Sprintf("ERROR (%s:%d): %s", ctx.FileName, ctx.Line, message)
}

type TokenType uint8

const (
	TypeInvalid TokenType = iota
	TypeNoOp
	TypePush
	TypePushPtr
	TypePushStr
	TypeGetStr
	TypePop
	TypeDup
	TypeCall
	TypeRet
	TypeEntrypoint
	TypePopStr
	TypeDupStr
	TypeInDupStr
	TypeSwapStr
	TypeInSwapStr
	TypeInDup
	TypeSwap
	TypeInSwap
	TypeAdd
	TypeSub
	TypeMul
	TypeDiv
	TypeMod
	TypeCmpe
	TypeCmpne
	TypeCmpg
	TypeCmpl
	TypeCmpge
	TypeCmple
	TypeJmp
	TypeZjmp
	TypeNzjmp
	TypePrint
	TypeNative
	TypeHalt
	TypeIntToStr
	TypeInt
	TypeFloat
	TypeChar
	TypeString
	TypeLabelDefinition
	TypeLabel
	TypeNull
)

type Token struct {
	Type      TokenType
	Text      string
	Line      int64
	Character int
	FileName  string
}
