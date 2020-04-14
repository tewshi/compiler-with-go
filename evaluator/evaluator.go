package evaluator

import (
	"fmt"
	"math"
	"monkey/ast"
	"monkey/object"
	"monkey/token"
	"monkey/utils"
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

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROROBJ
	}
	return false
}

// Eval returns the evaluated node as an object
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalProgram(node, env)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)

		// if the eval value is an identifier, then fetch out the value from the identifier
		if val.Type() == object.IDENTIFIEROBJ {
			val = val.(*object.Identifier).Value
		}
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}

		// Expressions
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) >= 1 {
			for _, arg := range args {
				if isError(arg) {
					return arg
				}
			}
		}

		return applyFunction(function, args)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		// stop propagation here if we encounter an error
		if isError(right) {
			return right
		}

		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		// stop propagation here if we encounter an error
		if isError(left) {
			return left
		}

		right := Eval(node.Right, env)
		// stop propagation here if we encounter an error
		if isError(right) {
			return right
		}

		modOps := []interface{}{
			token.PLUSEQ,
			token.MINUSEQ,
			token.SLASHEQ,
			token.ASTERISKEQ,
		}

		var val object.Object = evalInfixExpression(node.Operator, left, right)

		if utils.InArray(node.Operator, modOps) {
			if val.Type() == object.ERROROBJ {
				return val
			}

			switch node.Left.(type) {
			case *ast.Identifier:
				env.Set(node.Left.String(), val)
			}
		}

		return val

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.DoubleLiteral:
		return &object.Double{Value: node.Value, Precision: node.Precision}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)

		if len(elements) >= 1 {
			for _, element := range elements {
				if isError(element) {
					return element
				}
			}
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		// stop propagation here if we encounter an error
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	}
	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Identifier:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURNVALUEOBJ || rt == object.ERROROBJ {
				return result
			}
		}
	}
	return result
}

func evalExpressions(exps ast.Expressions, env *object.Environment) object.Objects {
	var result object.Objects
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return object.Objects{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for keyNode, valueNode := range node.Pairs {
		key := Eval(keyNode, env)
		if isError(key) {
			return key
		}
		if key.Type() == object.IDENTIFIEROBJ {
			key = key.(*object.Identifier).Value
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("unusable as hash key: %s", key.Type())
		}
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: value}
	}
	return &object.Hash{Pairs: pairs}
}

func evalIndexExpression(left object.Object, index object.Object) object.Object {
	var l object.Object = left
	var idx object.Object = index
	if left.Type() == object.IDENTIFIEROBJ {
		l = left.(*object.Identifier).Value
	}
	if index.Type() == object.IDENTIFIEROBJ {
		idx = index.(*object.Identifier).Value
	}
	switch {
	case l.Type() == object.ARRAYOBJ && idx.Type() == object.INTEGEROBJ:
		return evalArrayIndexExpression(l, idx)
	case l.Type() == object.HASHOBJ:
		return evalHashIndexExpression(l, idx)
	default:
		return newError("index operator not supported: %s", l.Type())
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)
	if idx < 0 || idx > max {
		return newError("array index out of bounds[0, %d]: %d", max, idx)
	}
	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index.Type())
	}
	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

func applyFunction(fn object.Object, args object.Objects) object.Object {
	switch fn := fn.(type) {
	case *object.Identifier:
		return applyFunction(fn.Value, args)
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args object.Objects) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {

	if val, ok := env.Get(node.Value); ok {
		if val.Type() == object.IDENTIFIEROBJ {
			val = val.(*object.Identifier).Value
		}
		return &object.Identifier{Name: node.Value, Value: val}
	}

	if builtin, ok := builtins[node.Value]; ok {
		return &object.Identifier{Name: node.Value, Value: builtin}
	}

	return newError("identifier not found: " + node.Value)
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	// stop propagation here if we encounter an error
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
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
	var r object.Object = right
	if right.Type() == object.IDENTIFIEROBJ {
		r = right.(*object.Identifier).Value
	}

	switch operator {
	case "!":
		return evalBangOperatorExpression(r)
	case "-":
		return evalMinusPrefixOperatorExpression(r)
	default:
		return newError("unknown operator: %s%s", operator, r.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	var l object.Object = left
	var r object.Object = right
	if left.Type() == object.IDENTIFIEROBJ {
		l = left.(*object.Identifier).Value
	}
	if right.Type() == object.IDENTIFIEROBJ {
		r = right.(*object.Identifier).Value
	}
	switch {
	case l.Type() == object.INTEGEROBJ && r.Type() == object.INTEGEROBJ:
		return evalIntegerInfixExpression(operator, l, r)
	case l.Type() == object.DOUBLEOBJ && r.Type() == object.DOUBLEOBJ:
		return evalDoubleInfixExpression(operator, l, r)
	case l.Type() == object.BOOLEANOBJ && r.Type() == object.BOOLEANOBJ:
		return evalBooleanInfixExpression(operator, l, r)
	case l.Type() == object.STRINGOBJ && r.Type() == object.STRINGOBJ:
		return evalStringInfixExpression(operator, l, r)
	case operator == token.POWER:
		switch {
		case l.Type() == object.INTEGEROBJ && r.Type() == object.DOUBLEOBJ:
			return evalPowerOperatorDoubleIntegerExpression(l, r)
		case l.Type() == object.DOUBLEOBJ && r.Type() == object.INTEGEROBJ:
			return evalPowerOperatorDoubleIntegerExpression(l, r)
		default:
			return newError("type mismatch: %s %s %s", l.Type(), operator, r.Type())
		}
	case l.Type() != r.Type():
		return newError("type mismatch: %s %s %s", l.Type(), operator, r.Type())
	case operator == token.EQ:
		return nativeBoolToBooleanObject(l == r)
	case operator == token.NOTEQ:
		return nativeBoolToBooleanObject(l != r)
	default:
		return newError("unknown operator: %s %s %s", l.Type(), operator, r.Type())
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

	// < <= > >=
	case token.LT:
		return evalLessThanOperatorIntegerExpression(left, right)
	case token.LTEQ:
		return evalLessThanEqualToOperatorIntegerExpression(left, right)
	case token.GT:
		return evalGreaterThanOperatorIntegerExpression(left, right)
	case token.GTEQ:
		return evalGreaterThanEqualToOperatorIntegerExpression(left, right)

	// == !=
	case token.EQ:
		return evalEqualToOperatorIntegerExpression(left, right)
	case token.NOTEQ:
		return evalNotEqualToOperatorIntegerExpression(left, right)

	// += -= *= /=
	case token.PLUSEQ:
		return evalPlusOperatorIntegerExpression(left, right)
	case token.MINUSEQ:
		return evalSubtractOperatorIntegerExpression(left, right)
	case token.ASTERISKEQ:
		return evalMultiplyOperatorIntegerExpression(left, right)
	case token.SLASHEQ:
		return evalDivideOperatorIntegerExpression(left, right)

	// ^
	case token.POWER:
		return evalPowerOperatorIntegerExpression(left, right)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalDoubleInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	if left.Type() != object.DOUBLEOBJ || right.Type() != object.DOUBLEOBJ {
		return newError("type mismatch: %s + %s", left.Type(), right.Type())
	}

	l := left.(*object.Double)
	r := right.(*object.Double)
	lvalue := l.Value
	rvalue := r.Value
	precision := utils.MaxInt(l.Precision, r.Precision)

	switch operator {
	// + - * /
	case token.PLUS:
		return &object.Double{Value: lvalue + rvalue, Precision: precision}
	case token.MINUS:
		return &object.Double{Value: lvalue - rvalue, Precision: precision}
	case token.ASTERISK:
		product := lvalue * rvalue
		precision := utils.MaxInt(precision, utils.Precision(fmt.Sprint(product)))
		return &object.Double{Value: product, Precision: precision}
	case token.SLASH:
		div := lvalue / rvalue
		precision := utils.MaxInt(precision, utils.Precision(fmt.Sprint(div)))
		return &object.Double{Value: div, Precision: precision}

	// < <= > >=
	case token.LT:
		return nativeBoolToBooleanObject(lvalue < rvalue)
	case token.LTEQ:
		return nativeBoolToBooleanObject(lvalue <= rvalue)
	case token.GT:
		return nativeBoolToBooleanObject(lvalue > rvalue)
	case token.GTEQ:
		return nativeBoolToBooleanObject(lvalue >= rvalue)

	// == !=
	case token.EQ:
		return nativeBoolToBooleanObject(lvalue == rvalue)
	case token.NOTEQ:
		return nativeBoolToBooleanObject(lvalue != rvalue)

	// ^
	case token.POWER:
		return evalPowerOperatorDoubleIntegerExpression(left, right)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPowerOperatorDoubleIntegerExpression(left object.Object, right object.Object) object.Object {

	if !((left.Type() != object.INTEGEROBJ || left.Type() != object.DOUBLEOBJ) && (right.Type() != object.INTEGEROBJ || right.Type() != object.DOUBLEOBJ)) {
		return newError("type mismatch: %s ^ %s", left.Type(), right.Type())
	}

	var lvalue float64
	var rvalue float64

	if left.Type() == object.DOUBLEOBJ {
		lvalue = left.(*object.Double).Value
	} else {
		lvalue = float64(left.(*object.Integer).Value)
	}

	if right.Type() == object.DOUBLEOBJ {
		rvalue = right.(*object.Double).Value
	} else {
		rvalue = float64(right.(*object.Integer).Value)
	}

	val := math.Pow(lvalue, rvalue)

	precision := utils.Precision(fmt.Sprint(val))

	return &object.Double{Value: val, Precision: precision}
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

func evalPowerOperatorIntegerExpression(left object.Object, right object.Object) object.Object {
	if left.Type() != object.INTEGEROBJ || right.Type() != object.INTEGEROBJ {
		return newError("type mismatch: %s / %s", left.Type(), right.Type())
	}

	lvalue := left.(*object.Integer).Value
	rvalue := right.(*object.Integer).Value

	val := int64(math.Pow(float64(lvalue), float64(rvalue)))

	return &object.Integer{Value: val}
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

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	lvalue := left.(*object.String).Value
	rvalue := right.(*object.String).Value
	switch operator {
	case token.NOTEQ:
		return nativeBoolToBooleanObject(lvalue != rvalue)
	case token.EQ:
		return nativeBoolToBooleanObject(lvalue == rvalue)
	case token.PLUS:
		return &object.String{Value: lvalue + rvalue}
	case token.PLUSEQ:
		return &object.String{Value: lvalue + rvalue}
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
