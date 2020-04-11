package object

import "fmt"

const (
	// INTEGEROBJ represents an integer object
	INTEGEROBJ = "INTEGER"
	// BOOLEANOBJ represents an boolean object
	BOOLEANOBJ = "BOOLEAN"
	// NULLOBJ represents an nil object
	NULLOBJ = "NULL"
	// NANOBJ represents an nil object
	NANOBJ = "NAN"
	// RETURNVALUEOBJ represents a return object
	RETURNVALUEOBJ = "RETURN_VALUE"
)

// Type represents the type of an object
type Type string

// Object wraps every value of the language
type Object interface {
	Type() Type
	Inspect() string
}

// Integer the int type
type Integer struct {
	Value int64
}

// Type returns the object type of this value
func (i *Integer) Type() Type { return INTEGEROBJ }

// Inspect inspector of the integer value
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean the bool type
type Boolean struct {
	Value bool
}

// Type returns the object type of this value
func (i *Boolean) Type() Type { return BOOLEANOBJ }

// Inspect inspector of the boolean value
func (i *Boolean) Inspect() string { return fmt.Sprintf("%t", i.Value) }

// Null the nil type
type Null struct{}

// Type returns the object type of this value
func (i *Null) Type() Type { return NULLOBJ }

// Inspect inspector of the Null value
func (i *Null) Inspect() string { return "null" }

// Nan represents not-a-number
type Nan struct{}

// Type returns the object type of this value
func (i *Nan) Type() Type { return NANOBJ }

// Inspect inspector of the Nan value
func (i *Nan) Inspect() string { return "NAN" }

// ReturnValue represents a return value
type ReturnValue struct {
	Value Object
}

// Type returns the object type of this value
func (rv *ReturnValue) Type() Type { return RETURNVALUEOBJ }

// Inspect inspector of the Nan value
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }
