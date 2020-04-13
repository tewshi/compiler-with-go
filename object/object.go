package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

const (
	// INTEGEROBJ represents an integer object
	INTEGEROBJ = "INTEGER"
	// STRINGOBJ represents an string object
	STRINGOBJ = "STRING"
	// BOOLEANOBJ represents an boolean object
	BOOLEANOBJ = "BOOLEAN"
	// ARRAYOBJ represents an array object
	ARRAYOBJ = "ARRAY"
	// NULLOBJ represents an nil object
	NULLOBJ = "NULL"
	// NANOBJ represents an nil object
	NANOBJ = "NAN"
	// RETURNVALUEOBJ represents a return object
	RETURNVALUEOBJ = "RETURN_VALUE"
	// ERROROBJ represents an error object
	ERROROBJ = "ERROR"
	// IDENTIFIEROBJ represents an identifier object
	IDENTIFIEROBJ = "IDENTIFIER"
	// FUNCTIONOBJ represents a function object
	FUNCTIONOBJ = "FUNCTION"
	// BUILTINOBJ represents a built-in function object
	BUILTINOBJ = "BUILTIN"
)

// Type represents the type of an object
type Type string

// Object wraps every value of the language
type Object interface {
	Type() Type
	Inspect() string
}

// Objects list of objects
type Objects []Object

// BuiltinFunction represents a built-in function
type BuiltinFunction func(args ...Object) Object

// Integer the int type
type Integer struct {
	Value int64
}

// Type returns the object type of this value
func (i *Integer) Type() Type { return INTEGEROBJ }

// Inspect returns a readable string of the integer value
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// String the int type
type String struct {
	Value string
}

// Type returns the object type of this value
func (s *String) Type() Type { return STRINGOBJ }

// Inspect returns a readable string of the String value
func (s *String) Inspect() string { return s.Value }

// Boolean the bool type
type Boolean struct {
	Value bool
}

// Type returns the object type of this value
func (i *Boolean) Type() Type { return BOOLEANOBJ }

// Inspect returns a readable string of the boolean value
func (i *Boolean) Inspect() string { return fmt.Sprintf("%t", i.Value) }

// Array the array data structure
type Array struct {
	Elements Objects
}

// Type returns the object type of this value
func (ao *Array) Type() Type { return ARRAYOBJ }

// Inspect returns a readable string of the array value
func (ao *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// Null the nil type
type Null struct{}

// Type returns the object type of this value
func (i *Null) Type() Type { return NULLOBJ }

// Inspect returns a readable string of the Null value
func (i *Null) Inspect() string { return "null" }

// Nan represents not-a-number
type Nan struct{}

// Type returns the object type of this value
func (i *Nan) Type() Type { return NANOBJ }

// Inspect returns a readable string of the Nan value
func (i *Nan) Inspect() string { return "NAN" }

// ReturnValue represents a return value
type ReturnValue struct {
	Value Object
}

// Type returns the object type of this value
func (rv *ReturnValue) Type() Type { return RETURNVALUEOBJ }

// Inspect returns a readable string of the return value
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error represents an error in our program
type Error struct {
	Message string
}

// Type returns the object type of this value
func (e *Error) Type() Type { return ERROROBJ }

// Inspect returns a readable string of the error
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Identifier the int type
type Identifier struct {
	Name  string
	Value Object
}

// Type returns the object type of this value
func (i *Identifier) Type() Type { return IDENTIFIEROBJ }

// Inspect returns a readable string of the Identifier value
func (i *Identifier) Inspect() string { return i.Value.Inspect() }

// Function represents a function in our program
type Function struct {
	Parameters ast.Identifiers
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type returns the object type of this value
func (f *Function) Type() Type { return FUNCTIONOBJ }

// Inspect returns a readable string of the function
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString(" (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	for _, s := range f.Body.Statements {
		out.WriteString(ast.TAB + s.String() + "\n")
	}
	out.WriteString("};")
	return out.String()
}

// Builtin represents a built-in function in our program
type Builtin struct {
	Fn BuiltinFunction
}

// Type returns the object type of this value
func (b *Builtin) Type() Type { return BUILTINOBJ }

// Inspect returns a readable string of the build-in function
func (b *Builtin) Inspect() string { return "builtin function" }
