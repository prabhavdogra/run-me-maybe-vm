package lexer

type TokenType uint8

const (
	TypeNoOp TokenType = iota
	TypePush
	TypePop
	TypeDup
	TypeIndup
	TypeSwap
	TypeInswap
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
	TypeHalt
)

type Token struct {
	Type      TokenType
	Text      string
	Line      int
	Character int
}
