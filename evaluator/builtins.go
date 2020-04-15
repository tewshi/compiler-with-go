package evaluator

import (
	"monkey/object"
	"strings"
)

var builtins = map[string]*object.Builtin{
	"len":   &object.Builtin{Fn: _len},
	"first": &object.Builtin{Fn: _first},
	"last":  &object.Builtin{Fn: _last},
	"rest":  &object.Builtin{Fn: _rest},
	"push":  &object.Builtin{Fn: _push},
	"puts":  &object.Builtin{Fn: _puts},
	"type":  &object.Builtin{Fn: _type},
}

func _len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	var arg object.Object = args[0]

	if arg.Type() == object.IDENTIFIEROBJ {
		arg = arg.(*object.Identifier).Value
	}

	switch arg.(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.(*object.String).Value))}
	case *object.Array:
		return &object.Integer{Value: int64(len(arg.(*object.Array).Elements))}
	default:
		return newError("argument to `len` not supported, got %s, want %s or %s",
			args[0].Type(), object.STRINGOBJ, object.ARRAYOBJ)
	}
}

func _first(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	var arg object.Object = args[0]
	if arg.Type() == object.IDENTIFIEROBJ {
		arg = arg.(*object.Identifier).Value
	}

	if arg.Type() != object.ARRAYOBJ {
		return newError("argument to `first` must be ARRAY, got %s",
			arg.Type())
	}

	elements := arg.(*object.Array).Elements
	if len(elements) > 0 {
		return elements[0]
	}
	return NULL
}

func _last(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	var arg object.Object = args[0]
	if arg.Type() == object.IDENTIFIEROBJ {
		arg = arg.(*object.Identifier).Value
	}

	if arg.Type() != object.ARRAYOBJ {
		return newError("argument to `last` must be ARRAY, got %s",
			arg.Type())
	}

	elements := arg.(*object.Array).Elements
	if length := len(elements); length > 0 {
		return elements[length-1]
	}
	return NULL
}

func _rest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}
	var arg object.Object = args[0]
	if arg.Type() == object.IDENTIFIEROBJ {
		arg = arg.(*object.Identifier).Value
	}

	if arg.Type() != object.ARRAYOBJ {
		return newError("argument to `last` must be ARRAY, got %s",
			arg.Type())
	}

	elements := arg.(*object.Array).Elements
	if length := len(elements); length > 0 {
		return &object.Array{Elements: elements[1:length]}
	}
	return NULL
}

func _push(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newError("wrong number of arguments. got=%d, want=2",
			len(args))
	}
	var arg object.Object = args[0]
	var val object.Object = args[1]
	if arg.Type() == object.IDENTIFIEROBJ {
		arg = arg.(*object.Identifier).Value
	}
	if val.Type() == object.IDENTIFIEROBJ {
		val = val.(*object.Identifier).Value
	}

	if arg.Type() != object.ARRAYOBJ {
		return newError("first argument to `push` must be ARRAY, got %s",
			arg.Type())
	}

	elements := arg.(*object.Array).Elements
	return &object.Array{Elements: append(elements, val)}
}

func _puts(args ...object.Object) object.Object {

	value := []string{}

	for _, arg := range args {
		value = append(value, arg.Inspect())
	}

	return &object.String{Value: strings.Join(value, "\n")}
}

func _type(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newError("wrong number of arguments. got=%d, want=1",
			len(args))
	}

	var arg object.Object = args[0]
	if arg.Type() == object.IDENTIFIEROBJ {
		arg = arg.(*object.Identifier).Value
	}

	return &object.String{Value: strings.ToLower(string(arg.Type()))}
}
