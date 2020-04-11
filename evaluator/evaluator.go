package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
)

var (
	// NULL holds a single null value for reuse
	NULL = &object.Null{}

	// NAN holds a single nan value for reuse
	NAN = &object.Nan{}

	// TRUE holds a single true value for reuse
	TRUE = &object.Boolean{Value: true}

	// FALSE holds a single false value for reuse
	FALSE = &object.Boolean{Value: false}
)

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// Eval returns the evaluated node as an object
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// Expressions
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.ReturnStatement:
		val := Eval(node.Value)
		return &object.ReturnValue{Value: val}
	}
	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURNVALUEOBJ || rt == object.ERROROBJ {
				return result
			}
		}
	}
	return result
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	case TRUE:
		return true
	default:
		return true
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGEROBJ && right.Type() == object.INTEGEROBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEANOBJ && right.Type() == object.BOOLEANOBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case operator == token.EQ:
		return nativeBoolToBooleanObject(left == right)
	case operator == token.NOTEQ:
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	// + - * /
	case token.PLUS:
		return evalPlusOperatorIntegerExpression(left, right)
	case token.MINUS:
		return evalSubtractOperatorIntegerExpression(left, right)
	case token.ASTERISK:
		return evalMultiplyOperatorIntegerExpression(left, right)
	case token.SLASH:
		return evalDivideOperatorIntegerExpression(left, right)
	// < <= > >= == !=
	case token.LT:
		return evalLessThanOperatorIntegerExpression(left, right)
	case token.LTEQ:
		return evalLessThanEqualToOperatorIntegerExpression(left, right)
	case token.GT:
		return evalGreaterThanOperatorIntegerExpression(left, right)
	case token.GTEQ:
		return evalGreaterThanEqualToOperatorIntegerExpression(left, right)
	case token.EQ:
		return evalEqualToOperatorIntegerExpression(left, right)
	case token.NOTEQ:
		return evalNotEqualToOperatorIntegerExpression(left, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalPlusOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s + %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return &object.Integer{Value: lvalue + rvalue}
}

func evalSubtractOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s - %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return &object.Integer{Value: lvalue - rvalue}
}

func evalMultiplyOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s * %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return &object.Integer{Value: lvalue * rvalue}
}

func evalDivideOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s / %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	if rvalue == 0 {
		return NAN
	}

	return &object.Integer{Value: lvalue / rvalue}
}

func evalLessThanOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s < %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return nativeBoolToBooleanObject(lvalue < rvalue)
}

func evalLessThanEqualToOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s <= %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return nativeBoolToBooleanObject(lvalue <= rvalue)
}

func evalGreaterThanOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s > %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return nativeBoolToBooleanObject(lvalue > rvalue)
}

func evalGreaterThanEqualToOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s >= %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return nativeBoolToBooleanObject(lvalue >= rvalue)
}

func evalEqualToOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s == %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return nativeBoolToBooleanObject(lvalue == rvalue)
}

func evalNotEqualToOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s != %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	return nativeBoolToBooleanObject(lvalue != rvalue)
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGEROBJ {
		return newError("unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	lvalue := left.(*object.Boolean).Value
	rvalue := right.(*object.Boolean).Value
	switch operator {
	case token.NOTEQ:
		return nativeBoolToBooleanObject(lvalue != rvalue)
	case token.EQ:
		return nativeBoolToBooleanObject(lvalue == rvalue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
