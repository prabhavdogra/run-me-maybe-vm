package main

import "fmt"

type LiteralType uint8

const (
	LiteralNone LiteralType = iota
	LiteralInt
	LiteralFloat
	LiteralChar
	LiteralString
)

type Literal struct {
	valueType   LiteralType
	valueInt    int64
	valueFloat  float64
	valueChar   rune
	valueString string
}

func (l *Literal) Type() LiteralType {
	return l.valueType
}

func IntLiteral(value int64) Literal {
	return Literal{
		valueType: LiteralInt,
		valueInt:  value,
	}
}

func FloatLiteral(value float64) Literal {
	return Literal{
		valueType:  LiteralFloat,
		valueFloat: value,
	}
}

func CharLiteral(value rune) Literal {
	return Literal{
		valueType: LiteralChar,
		valueChar: value,
	}
}

func StringLiteral(value string) Literal {
	return Literal{
		valueType:   LiteralString,
		valueString: value,
	}
}

func (l Literal) String() string {
	if l.Type() == LiteralInt {
		return fmt.Sprintf("INT %d", l.valueInt)
	} else if l.Type() == LiteralFloat {
		return fmt.Sprintf("FLOAT %f", l.valueFloat)
	} else if l.Type() == LiteralChar {
		return fmt.Sprintf("CHAR %c", l.valueChar)
	} else if l.Type() == LiteralString {
		return l.valueString
	}
	return "NONE"
}

func (l Literal) Print() {
	fmt.Println(l.String())
}

func (l Literal) Equal(other Literal) bool {
	if l.Type() != other.Type() {
		return false
	}
	switch l.Type() {
	case LiteralInt:
		return l.valueInt == other.valueInt
	case LiteralFloat:
		return l.valueFloat == other.valueFloat
	case LiteralChar:
		return l.valueChar == other.valueChar
	case LiteralString:
		return l.valueString == other.valueString
	default:
		return false
	}
}

func (l Literal) Greater(other Literal) bool {
	if l.Type() != other.Type() {
		panic("ERROR: \"greater\" comparison requires operands of same type")
	}
	switch l.Type() {
	case LiteralInt:
		return l.valueInt > other.valueInt
	case LiteralFloat:
		return l.valueFloat > other.valueFloat
	default:
		panic("ERROR: \"greater\" comparison not supported for this type")
	}
}

func (l Literal) Less(other Literal) bool {
	if l.Type() != other.Type() {
		panic("ERROR: \"less\" comparison requires operands of same type")
	}
	switch l.Type() {
	case LiteralInt:
		return l.valueInt < other.valueInt
	case LiteralFloat:
		return l.valueFloat < other.valueFloat
	default:
		panic("ERROR: \"less\" comparison not supported for this type")
	}
}

func (l Literal) GreaterOrEqual(other Literal) bool {
	if l.Type() != other.Type() {
		panic("ERROR:\"greater or equal\" comparison requires operands of same type")
	}
	switch l.Type() {
	case LiteralInt:
		return l.valueInt >= other.valueInt
	case LiteralFloat:
		return l.valueFloat >= other.valueFloat
	default:
		panic("ERROR:\"greater or equal\" comparison not supported for this type")
	}
}

func (l Literal) LessOrEqual(other Literal) bool {
	if l.Type() != other.Type() {
		panic("ERROR: \"less or equal\" comparison requires operands of same type")
	}
	switch l.Type() {
	case LiteralInt:
		return l.valueInt <= other.valueInt
	case LiteralFloat:
		return l.valueFloat <= other.valueFloat
	default:
		panic("ERROR: \"less or equal\" comparison not supported for this type")
	}
}

func (l Literal) Add(other Literal) Literal {
	if l.Type() != other.Type() {
		panic("ERROR: \"add\" requires operands of same type")
	}
	switch l.Type() {
	case LiteralInt:
		return IntLiteral(l.valueInt + other.valueInt)
	case LiteralFloat:
		return FloatLiteral(l.valueFloat + other.valueFloat)
	default:
		panic("ERROR: \"add\" not supported for this type")
	}
}

func (l Literal) Sub(other Literal) Literal {
	if l.Type() != other.Type() {
		panic("ERROR: \"sub\" requires operands of same type")
	}
	switch l.Type() {
	case LiteralInt:
		return IntLiteral(l.valueInt - other.valueInt)
	case LiteralFloat:
		return FloatLiteral(l.valueFloat - other.valueFloat)
	default:
		panic("ERROR: \"sub\" not supported for this type")
	}
}

func (l Literal) Mul(other Literal) Literal {
	if l.Type() != other.Type() {
		panic("ERROR: \"mul\" requires operands of same type")
	}
	switch l.Type() {
	case LiteralInt:
		return IntLiteral(l.valueInt * other.valueInt)
	case LiteralFloat:
		return FloatLiteral(l.valueFloat * other.valueFloat)
	default:
		panic("ERROR: \"mul\" not supported for this type")
	}
}

func (l Literal) Div(other Literal) Literal {
	if l.Type() != other.Type() {
		panic("ERROR: \"div\" requires operands of same type")
	}
	switch l.Type() {
	case LiteralInt:
		if other.valueInt == 0 {
			panic("ERROR: division by zero")
		}
		return IntLiteral(l.valueInt / other.valueInt)
	case LiteralFloat:
		if other.valueFloat == 0.0 {
			panic("ERROR: division by zero")
		}
		return FloatLiteral(l.valueFloat / other.valueFloat)
	default:
		panic("ERROR: \"div\" not supported for this type")
	}
}
