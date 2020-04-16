package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strconv"
	"strings"
)

const (
	// INTEGEROBJ represents an integer object
	INTEGEROBJ = "INTEGER"
	// DOUBLEOBJ represents a double object
	DOUBLEOBJ = "DOUBLE"
	// STRINGOBJ represents an string object
	STRINGOBJ = "STRING"
	// BOOLEANOBJ represents an boolean object
	BOOLEANOBJ = "BOOLEAN"
	// ARRAYOBJ represents an array object
	ARRAYOBJ = "ARRAY"
	// HASHOBJ represents an hash object
	HASHOBJ = "HASH"
	// NULLOBJ represents an nil object
	NULLOBJ = "NULL"
	// NANOBJ represents an nil object
	NANOBJ = "NAN"
	// COMMENTOBJ represents an nil object
	COMMENTOBJ = "COMMENT"
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

// Hashable represents a hashable object
type Hashable interface {
	HashKey() HashKey
}

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

// HashKey generates a hash for an integer key
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Double the int type
type Double struct {
	Value     float64
	Precision int // the double's precision
}

// Type returns the object type of this value
func (db *Double) Type() Type { return DOUBLEOBJ }

// Inspect returns a readable string of the double value
func (db *Double) Inspect() string { return strconv.FormatFloat(db.Value, 'f', db.Precision, 64) }

// HashKey generates a hash for an double key
func (db *Double) HashKey() HashKey {
	return HashKey{Type: db.Type(), Value: uint64(db.Value)}
}

// String the int type
type String struct {
	Value string
}

// Type returns the object type of this value
func (s *String) Type() Type { return STRINGOBJ }

// Inspect returns a readable string of the String value
func (s *String) Inspect() string { return s.Value }

// HashKey generates a hash for a string key
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Boolean the bool type
type Boolean struct {
	Value bool
}

// Type returns the object type of this value
func (b *Boolean) Type() Type { return BOOLEANOBJ }

// Inspect returns a readable string of the boolean value
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// HashKey generates a hash for a boolean key
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

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

// HashKey the has key object
type HashKey struct {
	Type  Type
	Value uint64
}

// HashPair represents a hash pair object... k:v pairs
type HashPair struct {
	Key   Object
	Value Object
}

// Hash represents a hash object... {k:v}
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type returns the object type of this value
func (h *Hash) Type() Type { return HASHOBJ }

// Inspect returns a readable string of the hash value
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
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

// Comment the comment type
type Comment struct {
	Value string
}

// Type returns the object type of this value
func (c *Comment) Type() Type { return COMMENTOBJ }

// Inspect returns a readable string of the comment value
func (c *Comment) Inspect() string { return c.Value }

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
