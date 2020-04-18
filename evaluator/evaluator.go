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
		return evalLetStatement(node, env)

	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)

	// Expressions
	case *ast.CallExpression:
		return evalFunctionCall(node, env)

	case *ast.Identifier:
		return evalIdentifier(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)

	case *ast.InfixExpression:
		return evalInfixExpression(node, env)

	case *ast.SuffixExpression:
		return evalSuffixExpression(node, env)

	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)

	case *ast.HashLiteral:
		return evalHashLiteral(node, env)

	case *ast.IndexExpression:
		return evalIndexExpression(node, env)

	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.DoubleLiteral:
		return &object.Double{Value: node.Value, Precision: node.Precision}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.NullLiteral:
		return NULL

		/* Comments never get here, lexer recognizes them, but the parser strips them off
		* case *ast.CommentLiteral:
		* // we return nil coz we dont evaluate comments
		* // return &object.Comment{Value: node.Value}
		* return nil
		 */
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

func evalLetStatement(stmt *ast.LetStatement, env *object.Environment) object.Object {
	val := Eval(stmt.Value, env)

	// if the eval value is an identifier, then fetch out the value from the identifier
	if val.Type() == object.IDENTIFIEROBJ {
		val = val.(*object.Identifier).Value
	}
	if isError(val) {
		return val
	}
	env.Set(stmt.Name.Value, val)
	return nil
}

func evalFunctionLiteral(fn *ast.FunctionLiteral, env *object.Environment) object.Object {
	params := fn.Parameters
	body := fn.Body
	return &object.Function{Parameters: params, Env: env, Body: body}
}

func evalFunctionCall(fn *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(fn.Function, env)
	if isError(function) {
		return function
	}

	args := evalExpressions(fn.Arguments, env)
	if len(args) >= 1 {
		for _, arg := range args {
			if isError(arg) {
				return arg
			}
		}
	}

	return applyFunction(function, args)
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

func evalPrefixExpression(pref *ast.PrefixExpression, env *object.Environment) object.Object {
	right := Eval(pref.Right, env)
	// stop propagation here if we encounter an error
	if isError(right) {
		return right
	}

	r := right
	isIdentifier := false
	var identifier string

	if right.Type() == object.IDENTIFIEROBJ {
		right := right.(*object.Identifier)
		r = right.Value
		isIdentifier = true
		identifier = right.Name
	}

	switch pref.Operator {
	case token.BANG:
		return evalBangOperatorExpression(r)
	case token.MINUS:
		switch {
		case r.Type() == object.INTEGEROBJ:
			return evalMinusPrefixOperatorExpression(r)
		case r.Type() == object.DOUBLEOBJ:
			return evalMinusPrefixOperatorExpression(r)
		default:
			return newError("unknown operator: %s%s", pref.Operator, r.Type())
		}
	case token.INCREMENT:

		var result object.Object

		switch {
		case isIdentifier && r.Type() == object.INTEGEROBJ:
			result = evalIncrementOperatorExpression(r)
		case isIdentifier && r.Type() == object.DOUBLEOBJ:
			result = evalIncrementOperatorExpression(r)
		default:
			result = newError("unknown operator: %s%s", pref.Operator, r.Type())
		}

		if result.Type() == object.ERROROBJ {
			return result
		}

		env.Set(identifier, result)

		return result

	case token.DECREMENT:

		var result object.Object

		switch {
		case isIdentifier && r.Type() == object.INTEGEROBJ:
			result = evalDecrementOperatorExpression(r)
		case isIdentifier && r.Type() == object.DOUBLEOBJ:
			result = evalDecrementOperatorExpression(r)
		default:
			result = newError("unknown operator: %s%s", pref.Operator, r.Type())
		}

		if result.Type() == object.ERROROBJ {
			return result
		}

		env.Set(identifier, result)

		return result

	default:
		return newError("unknown operator: %s%s", pref.Operator, r.Type())
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

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.DOUBLEOBJ:
		value := right.(*object.Double).Value
		return &object.Double{Value: -value}
	default:
		value := right.(*object.Integer).Value
		return &object.Integer{Value: -value}
	}
}

func evalIncrementOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.DOUBLEOBJ:
		r := right.(*object.Double)
		value := r.Value
		return evalPlusOperatorDoubleExpression(value, 1, r.Precision)
	default:
		value := right.(*object.Integer).Value
		return evalPlusOperatorIntegerExpression(value, 1)
	}
}

func evalDecrementOperatorExpression(right object.Object) object.Object {
	switch right.Type() {
	case object.DOUBLEOBJ:
		r := right.(*object.Double)
		value := r.Value
		return evalSubtractOperatorDoubleExpression(value, 1, r.Precision)
	default:
		value := right.(*object.Integer).Value
		return evalSubtractOperatorIntegerExpression(value, 1)
	}
}

func evalInfixExpression(inf *ast.InfixExpression, env *object.Environment) object.Object {
	left := Eval(inf.Left, env)
	// stop propagation here if we encounter an error
	if isError(left) {
		return left
	}

	right := Eval(inf.Right, env)
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

	var val object.Object = evalInfixExpressionByType(inf.Operator, left, right)

	if utils.InArray(inf.Operator, modOps) {
		if val.Type() == object.ERROROBJ {
			return val
		}

		switch inf.Left.(type) {
		case *ast.Identifier:
			env.Set(inf.Left.String(), val)
		}

		return nil
	}

	return val
}

func evalInfixExpressionByType(operator string, left object.Object, right object.Object) object.Object {
	var l object.Object = left
	var r object.Object = right
	if left.Type() == object.IDENTIFIEROBJ {
		l = left.(*object.Identifier).Value
	}
	if right.Type() == object.IDENTIFIEROBJ {
		r = right.(*object.Identifier).Value
	}
	switch {
	case operator == token.NOTNULLOR:
		return evalNullOrOperatorExpression(l, r)

	case operator == token.POWER:
		switch {
		case l.Type() == object.INTEGEROBJ && r.Type() == object.DOUBLEOBJ:
			return evalPowerOperatorDoubleIntegerExpression(l, r)
		case l.Type() == object.DOUBLEOBJ && r.Type() == object.INTEGEROBJ:
			return evalPowerOperatorDoubleIntegerExpression(l, r)
		case l.Type() == object.INTEGEROBJ && r.Type() == object.INTEGEROBJ:
			return evalIntegerInfixExpression(operator, l, r)
		case l.Type() == object.DOUBLEOBJ && r.Type() == object.DOUBLEOBJ:
			return evalDoubleInfixExpression(operator, l, r)
		default:
			return newError("type mismatch: %s %s %s", l.Type(), operator, r.Type())
		}

	case operator == token.MODULUS:
		switch {
		case l.Type() == object.INTEGEROBJ && r.Type() == object.DOUBLEOBJ:
			return evalModulusOperatorDoubleIntegerExpression(l, r)
		case l.Type() == object.DOUBLEOBJ && r.Type() == object.INTEGEROBJ:
			return evalModulusOperatorDoubleIntegerExpression(l, r)
		case l.Type() == object.INTEGEROBJ && r.Type() == object.INTEGEROBJ:
			return evalIntegerInfixExpression(operator, l, r)
		case l.Type() == object.DOUBLEOBJ && r.Type() == object.DOUBLEOBJ:
			return evalDoubleInfixExpression(operator, l, r)
		default:
			return newError("type mismatch: %s %s %s", l.Type(), operator, r.Type())
		}

	case l.Type() == object.INTEGEROBJ && r.Type() == object.INTEGEROBJ:
		return evalIntegerInfixExpression(operator, l, r)
	case l.Type() == object.DOUBLEOBJ && r.Type() == object.DOUBLEOBJ:
		return evalDoubleInfixExpression(operator, l, r)
	case l.Type() == object.BOOLEANOBJ && r.Type() == object.BOOLEANOBJ:
		return evalBooleanInfixExpression(operator, l, r)
	case l.Type() == object.STRINGOBJ && r.Type() == object.STRINGOBJ:
		return evalStringInfixExpression(operator, l, r)

	// + += - -= * *= / /=
	// < <= > >=
	// == !=
	case operator == token.PLUS, operator == token.PLUSEQ, operator == token.MINUS,
		operator == token.MINUSEQ, operator == token.ASTERISK, operator == token.ASTERISKEQ,
		operator == token.SLASH, operator == token.SLASHEQ,
		operator == token.LT, operator == token.LTEQ, operator == token.GT,
		operator == token.GTEQ, operator == token.EQ, operator == token.NOTEQ:
		switch {
		case l.Type() == object.INTEGEROBJ && r.Type() == object.DOUBLEOBJ:
			dl := &object.Double{Value: float64(l.(*object.Integer).Value), Precision: 0}
			return evalDoubleInfixExpression(operator, dl, r)
		case l.Type() == object.DOUBLEOBJ && r.Type() == object.INTEGEROBJ:
			dr := &object.Double{Value: float64(r.(*object.Integer).Value), Precision: 0}
			return evalDoubleInfixExpression(operator, l, dr)
		case operator == token.EQ:
			return nativeBoolToBooleanObject(l == r)
		case operator == token.NOTEQ:
			return nativeBoolToBooleanObject(l != r)
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

func evalSuffixExpression(suf *ast.SuffixExpression, env *object.Environment) object.Object {
	left := Eval(suf.Left, env)
	// stop propagation here if we encounter an error
	if isError(left) {
		return left
	}

	l := left
	operator := suf.Operator
	isIdentifier := false
	var identifier string

	if left.Type() == object.IDENTIFIEROBJ {
		left := left.(*object.Identifier)
		l = left.Value
		isIdentifier = true
		identifier = left.Name
	}

	switch operator {

	case token.INCREMENT:

		var result object.Object

		switch {
		case isIdentifier && l.Type() == object.INTEGEROBJ:
			result = evalIncrementOperatorExpression(l)
		case isIdentifier && l.Type() == object.DOUBLEOBJ:
			result = evalIncrementOperatorExpression(l)
		default:
			result = newError("unknown operator: %s%s", l.Type(), operator)
		}

		if result.Type() == object.ERROROBJ {
			return result
		}

		env.Set(identifier, result)

		return l

	case token.DECREMENT:

		var result object.Object

		switch {
		case isIdentifier && l.Type() == object.INTEGEROBJ:
			result = evalDecrementOperatorExpression(l)
		case isIdentifier && l.Type() == object.DOUBLEOBJ:
			result = evalDecrementOperatorExpression(l)
		default:
			result = newError("unknown operator: %s%s", l.Type(), operator)
		}

		if result.Type() == object.ERROROBJ {
			return result
		}

		env.Set(identifier, result)

		return l

	default:
		return newError("unknown operator: %s %s %s", l.Type(), operator, l.Type())
	}
}

func evalArrayLiteral(al *ast.ArrayLiteral, env *object.Environment) object.Object {
	elements := evalExpressions(al.Elements, env)

	if len(elements) >= 1 {
		for _, element := range elements {
			if isError(element) {
				return element
			}
		}
	}
	return &object.Array{Elements: elements}
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

func evalIndexExpression(ie *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(ie.Left, env)
	if isError(left) {
		return left
	}
	index := Eval(ie.Index, env)
	if isError(index) {
		return index
	}

	if left.Type() == object.IDENTIFIEROBJ {
		left = left.(*object.Identifier).Value
	}
	if index.Type() == object.IDENTIFIEROBJ {
		index = index.(*object.Identifier).Value
	}
	switch {
	case left.Type() == object.ARRAYOBJ && index.Type() == object.INTEGEROBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASHOBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s", left.Type())
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

func evalReturnStatement(rs *ast.ReturnStatement, env *object.Environment) object.Object {
	val := Eval(rs.Value, env)
	// stop propagation here if we encounter an error
	if isError(val) {
		return val
	}

	return &object.ReturnValue{Value: val}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	l := left.(*object.Integer)
	r := right.(*object.Integer)
	lvalue := l.Value
	rvalue := r.Value

	switch operator {
	// + += - -= * *= / /=
	case token.PLUS, token.PLUSEQ:
		return evalPlusOperatorIntegerExpression(lvalue, rvalue)
	case token.MINUS, token.MINUSEQ:
		return evalSubtractOperatorIntegerExpression(lvalue, rvalue)
	case token.ASTERISK, token.ASTERISKEQ:
		return evalMultiplyOperatorIntegerExpression(lvalue, rvalue)
	case token.SLASH, token.SLASHEQ:
		return evalDivideOperatorIntegerExpression(lvalue, rvalue)

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
		val := int64(math.Pow(float64(lvalue), float64(rvalue)))
		return &object.Integer{Value: val}
	case token.MODULUS:
		val := int64(math.Mod(float64(lvalue), float64(rvalue)))
		return &object.Integer{Value: val}

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalDoubleInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	l := left.(*object.Double)
	r := right.(*object.Double)
	lvalue := l.Value
	rvalue := r.Value
	precision := utils.MaxInt(l.Precision, r.Precision)

	switch operator {
	// + += - -= * *= / /=
	case token.PLUS, token.PLUSEQ:
		return evalPlusOperatorDoubleExpression(lvalue, rvalue, precision)
	case token.MINUS, token.MINUSEQ:
		return evalSubtractOperatorDoubleExpression(lvalue, rvalue, precision)
	case token.ASTERISK, token.ASTERISKEQ:
		return evalMultiplyOperatorDoubleExpression(lvalue, rvalue, precision)
	case token.SLASH, token.SLASHEQ:
		return evalDivideOperatorDoubleExpression(lvalue, rvalue, precision)

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

	// ^ %
	case token.POWER:
		return evalPowerOperatorDoubleIntegerExpression(left, right)
	case token.MODULUS:
		return evalModulusOperatorDoubleIntegerExpression(left, right)

	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPowerOperatorDoubleIntegerExpression(left object.Object, right object.Object) object.Object {
	if !((left.Type() == object.INTEGEROBJ || left.Type() == object.DOUBLEOBJ) && (right.Type() == object.INTEGEROBJ || right.Type() == object.DOUBLEOBJ)) {
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

func evalModulusOperatorDoubleIntegerExpression(left object.Object, right object.Object) object.Object {

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

	val := math.Mod(lvalue, rvalue)

	precision := utils.Precision(fmt.Sprint(val))

	return &object.Double{Value: val, Precision: precision}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	lvalue := left.(*object.Boolean).Value
	rvalue := right.(*object.Boolean).Value

	switch operator {
	// == !=
	case token.NOTEQ:
		return nativeBoolToBooleanObject(lvalue != rvalue)
	case token.EQ:
		return nativeBoolToBooleanObject(lvalue == rvalue)
	// && ||
	case token.AND:
		return nativeBoolToBooleanObject(lvalue && rvalue)
	case token.OR:
		return nativeBoolToBooleanObject(lvalue || rvalue)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPlusOperatorIntegerExpression(lvalue int64, rvalue int64) object.Object {
	return &object.Integer{Value: lvalue + rvalue}
}

func evalSubtractOperatorIntegerExpression(lvalue int64, rvalue int64) object.Object {
	return &object.Integer{Value: lvalue - rvalue}
}

func evalMultiplyOperatorIntegerExpression(lvalue int64, rvalue int64) object.Object {
	return &object.Integer{Value: lvalue * rvalue}
}

func evalDivideOperatorIntegerExpression(lvalue int64, rvalue int64) object.Object {
	return &object.Integer{Value: lvalue / rvalue}
}

func evalPlusOperatorDoubleExpression(lvalue float64, rvalue float64, precision int) object.Object {
	return &object.Double{Value: lvalue + rvalue, Precision: precision}
}

func evalSubtractOperatorDoubleExpression(lvalue float64, rvalue float64, precision int) object.Object {
	return &object.Double{Value: lvalue - rvalue, Precision: precision}
}

func evalMultiplyOperatorDoubleExpression(lvalue float64, rvalue float64, precision int) object.Object {
	product := lvalue * rvalue
	prec := utils.MaxInt(precision, utils.Precision(fmt.Sprint(product)))
	return &object.Double{Value: product, Precision: prec}
}

func evalDivideOperatorDoubleExpression(lvalue float64, rvalue float64, precision int) object.Object {
	div := lvalue / rvalue
	prec := utils.MaxInt(precision, utils.Precision(fmt.Sprint(div)))
	return &object.Double{Value: div, Precision: prec}
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

func evalNullOrOperatorExpression(left object.Object, right object.Object) object.Object {
	if left == NULL {
		return right
	}
	return left
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
