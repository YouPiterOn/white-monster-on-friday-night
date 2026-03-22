package compiler

import "fmt"

type PrimitiveType int

const (
	TYPE_ANY PrimitiveType = iota
	TYPE_UNKNOWN
	TYPE_INT
	TYPE_FLOAT
	TYPE_STRING
	TYPE_BOOL
	TYPE_NULL
	TYPE_VOID
	TYPE_ARRAY
	TYPE_OBJECT
	TYPE_FUNCTION
)

func (p PrimitiveType) String() string {
	return [...]string{
		"any",
		"unknown",
		"int",
		"float",
		"string",
		"bool",
		"null",
		"void",
		"array",
		"object",
		"function",
	}[p]
}

type FuncSignature struct {
	CallArgs   []*Type
	ReturnType *Type
	Vararg     bool
}

func (self *FuncSignature) IsSubsignatureOf(other *FuncSignature) bool {
	if len(self.CallArgs) != len(other.CallArgs) {
		return false
	}
	if self.Vararg != other.Vararg {
		return false
	}
	for i := range self.CallArgs {
		if !self.CallArgs[i].IsSubtypeOf(other.CallArgs[i]) {
			return false
		}
	}
	return self.ReturnType.IsSubtypeOf(other.ReturnType)
}

type Type struct {
	PrimitiveType PrimitiveType
	Proto         *Type
	ElementType   *Type
	Properties    map[string]*Type
	FuncSignature *FuncSignature
}

func (self *Type) IsSubtypeOf(other *Type) bool {
	if self.PrimitiveType == TYPE_ANY || self.PrimitiveType == TYPE_UNKNOWN || other.PrimitiveType == TYPE_ANY {
		return true
	}
	if other.PrimitiveType == TYPE_UNKNOWN {
		return false
	}
	if self.Proto != nil {
		if other.Proto == nil {
			return false
		}
		return self.Proto.IsSubtypeOf(other.Proto)
	}
	if self.PrimitiveType == TYPE_FUNCTION {
		if other.PrimitiveType != TYPE_FUNCTION {
			return false
		}
		return self.FuncSignature.IsSubsignatureOf(other.FuncSignature)
	}
	if self.PrimitiveType == TYPE_ARRAY {
		if other.PrimitiveType != TYPE_ARRAY {
			return false
		}
		return self.ElementType.IsSubtypeOf(other.ElementType)
	}
	if self.PrimitiveType == TYPE_OBJECT {
		if other.PrimitiveType != TYPE_OBJECT {
			return false
		}
		for key, value := range self.Properties {
			if otherValue, ok := other.Properties[key]; !ok || !value.IsSubtypeOf(otherValue) {
				return false
			}
		}
		return true
	}

	return self.PrimitiveType == other.PrimitiveType
}

func TypeInt() *Type {
	return &Type{Name: TYPE_INT.String()}
}

func TypeFloat() *Type {
	return &Type{Name: TYPE_FLOAT.String()}
}

func TypeString() *Type {
	return &Type{Name: TYPE_STRING.String()}
}

func TypeBool() *Type {
	return &Type{Name: TYPE_BOOL.String()}
}

func TypeNull() *Type {
	return &Type{Name: TYPE_NULL.String()}
}

func TypeVoid() *Type {
	return &Type{Name: TYPE_VOID.String()}
}

func TypeArrayOf(elementType *Type) *Type {
	return &Type{Name: fmt.Sprintf("%s[]", elementType.Name), ElementType: elementType}
}
