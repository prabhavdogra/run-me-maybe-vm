package token

type TokenType uint8

const (
	TypeInvalid TokenType = iota
	TypeNoOp
	TypePush
	TypePop
	TypeDup
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
	TypeHalt
	TypeInt
	TypeLabelDefinition
	TypeLabel
)

type Token struct {
	Type      TokenType
	Text      string
	Line      int64
	Character int
}
