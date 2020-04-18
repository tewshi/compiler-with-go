package evaluator

import (
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"5", 5},
		{"10", 10},

		{"-5", -5},
		{"-10", -10},

		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},

		{"5 ^ 2", 25},
		{"5 ^ -1", 0},
		{"5 ^ 1 + 5", 10},
		{"5 * 5 ^ 0", 5},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},

		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},

		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},

		{"let x = 1;x += 2; x == 3", true},
		{"let x = 1;x += x; x == 2", true},
		{"let x = 1;x += 2; x *= x; x == 9", true},
		{"let x = 1;x -= 2; x == -1", true},
		{"let x = 1;x -= x; x == 0", true},
		{"let x = 1;x -= 2; x *= x; x == 1", true},
		{"let x = 1;x *= 2; x == 2", true},
		{"let x = 1;x *= x; x == 1", true},
		{"let x = 1;x *= 2; x *= x; x == 4", true},
		{"let x = 4;x /= 2; x == 2", true},
		{"let x = 4;x *= x; x == 16", true},
		{"let x = 4;x *= 2; x *= x; x == 64", true},

		{"let x = 1.0;x += 2.0; x == 3.0", true},
		{"let x = 1.0;x += x; x == 2.0", true},
		{"let x = 1.0;x += 2.0; x *= x; x == 9.0", true},
		{"let x = 1.0;x -= 2.0; x == -1.0", true},
		{"let x = 1.0;x -= x; x == 0.0", true},
		{"let x = 1.0;x -= 2.0; x *= x; x == 1.0", true},
		{"let x = 1.0;x *= 2.0; x == 2.0", true},
		{"let x = 1.0;x *= x; x == 1.0", true},
		{"let x = 1.0;x *= 2.0; x *= x; x == 4.0", true},
		{"let x = 4.0;x /= 2.0; x == 2.0", true},
		{"let x = 4.0;x *= x; x == 16.0", true},
		{"let x = 4.0;x *= 2.0; x *= x; x == 64.0", true},

		{"let x = 4.0;x ^ 2.0 == 16.0", true},
		{"let x = 4.0;x ^ 2 == 16.0", true},
		{"let x = 4;x ^ 2.0 == 16.0", true},
		{"let x = 4;4.0 + x ^ 2.0 == 4.0 ^ 2.0 + 4.0", true},

		{"let x = 4.0;x % 2.0 == 0.0", true},
		{"2 % 4 * 5^2 - 2 / 4 == 50", true},
		{"2 % 4 * 5^2 - 2 / 4.0 == 49.5", true},
		{"2 % 4 * 5^2 - 2.0 / 4 == 49.5", true},
		{"2 % 4 * 5^2.0 - 2 / 4 == 50", true},
		{"2 % 4 * 5.0^2 - 2 / 4 == 50", true},
		{"2 % 4.0 * 5^2 - 2 / 4 == 50", true},
		{"2.0 % 4 * 5^2 - 2 / 4 == 50", true},
		{"2.0 % 4.0 * 5.0^2.0 - 2.0 / 4.0 == 49.5", true},

		{"let x = 1.0;x++; x == 2.0", true},
		{"let x = 1.0;++x; x == 2.0", true},
		{"let x = 1.0;x++ == 1.0", true},
		{"let x = 1.0;++x == 2.0", true},
		{"let x = 1.0;x--; x == 0.0", true},
		{"let x = 1.0;--x; x == 0.0", true},
		{"let x = 1.0;x-- == 1.0", true},
		{"let x = 1.0;--x == 0.0", true},

		{"true || 2 % 4 * 5^2 >= 2 % 4 * 5^0", true},
		{"true && true", true},
		{"true && !true", false},
		{"true && false", false},
		{"!true && true", false},
		{"true || true", true},
		{"true || !true", true},
		{"true || false", true},
		{"!true || true", true},

		{"!nil", true},
		{"!!nil", false},
		{"nil == nil", true},
		{"nil != nil", false},
		{"false ?? false", false},
		{"true ?? true", true},
		{"false ?? true", false},
		{"true ?? false", true},
		{"let x = true; x ?? false", true},
		{"let x = nil; x ?? true", true},
		{"let x = nil; x ?? false", false},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"let x = 1; if (x == 1) { 10 }", 10},
		{"let x = 1; if (x++ == 1) { 10 }", 10},
		{"let x = 1; if (++x == 2) { 10 }", 10},
		{"let x = 1; if (x-- == 1) { 10 }", 10},
		{"let x = 1; if (--x == 0) { 10 }", 10},
		{"let x = 1; x++; if (x == 2) { 10 }", 10},
		{"let x = 1; ++x; if (x == 2) { 10 }", 10},
		{"let x = 1; x--; if (x == 0) { 10 }", 10},
		{"let x = 1; --x; if (x == 0) { 10 }", 10},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},

		{`if (10 > 1) {
			if (10 > 1) {
				return 10;
			}
			
			return 1;
		}`, 10},

		{"let addTwo = fn(x) { x + 2; }; addTwo(2);", 4},
		{"let multiply = fn(x, y) { x * y }; multiply(50 / 2, 1 * 2);", 50},
		{"let pow = fn(x) { x * x }; pow(5);", 25},
		{"let double = fn(x) { x * 2 }; double(5);", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return true + false;
				};
				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},

		{
			"foobar",
			"identifier not found: foobar",
		},

		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},

		{
			"[1, 2, 3][3]",
			"array index out of bounds[0, 2]: 3",
		},
		{
			"[1, 2, 3][-1]",
			"array index out of bounds[0, 2]: -1",
		},

		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},

		{
			"5 ^ true;",
			"type mismatch: INTEGER ^ BOOLEAN",
		},
		{
			"5 ^ \"hello\";",
			"type mismatch: INTEGER ^ STRING",
		},
		{
			`let x = 5; x += "Hello";`,
			"type mismatch: INTEGER += STRING",
		},
		{
			`let x = "5"; ++x;`,
			"unknown operator: ++STRING",
		},
		{
			`let x = "5"; --x;`,
			"unknown operator: --STRING",
		},
		{
			`let x = "5"; x++;`,
			"unknown operator: STRING++",
		},
		{
			`let x = "5"; x--;`,
			"unknown operator: STRING--",
		},
		{
			`let x = true; ++x;`,
			"unknown operator: ++BOOLEAN",
		},
		{
			`let x = true; --x;`,
			"unknown operator: --BOOLEAN",
		},
		{
			`let x = true; x++;`,
			"unknown operator: BOOLEAN++",
		},
		{
			`let x = true; x--;`,
			"unknown operator: BOOLEAN--",
		},
		{
			`let x = []; ++x;`,
			"unknown operator: ++ARRAY",
		},
		{
			`let x = []; --x;`,
			"unknown operator: --ARRAY",
		},
		{
			`let x = []; x++;`,
			"unknown operator: ARRAY++",
		},
		{
			`let x = []; x--;`,
			"unknown operator: ARRAY--",
		},
		{
			`let x = {}; ++x;`,
			"unknown operator: ++HASH",
		},
		{
			`let x = {}; --x;`,
			"unknown operator: --HASH",
		},
		{
			`let x = {}; x++;`,
			"unknown operator: HASH++",
		},
		{
			`let x = {}; x--;`,
			"unknown operator: HASH--",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}
		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}
	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringComparison(t *testing.T) {

	tests := []struct {
		input    string
		expected bool
	}{
		{`"x" == "x"`, true},
		{`"x" != "x"`, false},
		{`let x = "12345"; x == "12345"`, true},
		{`let x = "Hello"; let y = " World"; x + y == "Hello World"`, true},
		{`let x = "Hello"; let y = " World"; x += y; x == "Hello World"`, true},
		{`let x = "Hello"; x == x`, true},
		{`let x = "Hello"; let y = x; x == y`, true},
		{`let x = "Hello"; let y = "Hello"; x == y`, true},
		{`let x = "Hello"; let y = "World"; x == y`, false},
		{`let x = "Hello"; let y = "World"; x != y`, true},
		{`let x = "Hello"; let y = "World"; let z = x != y; z`, true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		str, ok := evaluated.(*object.Boolean)
		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
		}
		if str.Value != tt.expected {
			t.Errorf("String has wrong value. got=%t", str.Value)
		}
	}

}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER, want STRING or ARRAY"},
		{`len(true)`, "argument to `len` not supported, got BOOLEAN, want STRING or ARRAY"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},

		{`let a = [1, 2, 3, 4]; rest(a)`,
			&object.Array{Elements: object.Objects{
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
			}},
		},
		{`let a = [1, 2, 3, 4]; rest(rest(a))`,
			&object.Array{Elements: object.Objects{
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
			}},
		},

		{`let a = [1, 2, 3, 4]; let b = push(a, 5); a;`,
			&object.Array{Elements: object.Objects{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
			}},
		},
		{`let a = [1, 2, 3, 4]; let b = push(a, 5); b;`,
			&object.Array{Elements: object.Objects{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
				&object.Integer{Value: 5},
			}},
		},

		{`
		let map = fn(arr, f) {
			let iter = fn(arr, accumulated) {
				if (len(arr) == 0) {
					accumulated
				} else {
					iter(rest(arr), push(accumulated, f(first(arr))));
				}
			};
			
			iter(arr, []);
		};
		let a = [1, 2, 3, 4];
		let double = fn(x) { x * 2 };
		map(a, double);`,
			&object.Array{Elements: object.Objects{
				&object.Integer{Value: 2},
				&object.Integer{Value: 4},
				&object.Integer{Value: 6},
				&object.Integer{Value: 8},
			}},
		},
		{`
		let reduce = fn(arr, initial, f) {
			let iter = fn(arr, result) {
				if (len(arr) == 0) {
					result
				} else {
					iter(rest(arr), f(result, first(arr)));
				}
			};
			iter(arr, initial);
		};
		let sum = fn(arr) {
			reduce(arr, 0, fn(initial, el) { initial + el });
		};
		sum([1, 2, 3, 4, 5]);`, 15,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case *object.Array:
			arr, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("object is not Array. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}

			if len(arr.Elements) != len(expected.Elements) {
				t.Errorf("wrong array length. expected=%d, got=%d",
					len(expected.Elements), len(arr.Elements))
			}
			for i, el := range arr.Elements {
				if el.Inspect() != expected.Elements[i].Inspect() {
					t.Errorf("wrong object at index %d. expected=%q, got=%q",
						i, expected.Elements[i].Inspect(), el.Inspect())
				}

			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}
	t.Log(result)
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}
	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}
	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}
	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`let x = nil; x ?? nil`,
			nil,
		},
		{
			`nil ?? nil`,
			nil,
		},
		{
			`(1+1) ?? nil`,
			2,
		},
		{
			`100 ?? nil`,
			100,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
		{
			`{1.0: 10}[1.0]`,
			10,
		},
		{
			`{1.0: 1}[1.0]`,
			1,
		},
		{
			`{2^2: 16}[4]`,
			16,
		},
		{
			`{-2^2: 16}[4]`,
			16,
		},
		{
			`{-2^3: 16}[-8]`,
			16,
		},
		{
			`{-(2^3): 16}[-8]`,
			16,
		},
		{
			`{-(2^3): 16}[-8] + 4`,
			20,
		},
		{
			`let people = [{"name": "Alice", "age": 24}, {"name": "Anna", "age": 28}]; let n = "name"; let z = 0; len(people[z][n])`,
			5,
		},
		{
			`let people = [{"name": "Alice", "age": 24}, {"name": "Anna", "age": 28}]; let n = "name"; let z = 0; people[1]["age"]`,
			28,
		},
		{
			`let getAge = fn(person) { person["age"]; }; getAge({"name": "Alice", "age": 24});`,
			24,
		},
		{
			`let getAge = fn(person) { person["age"]; }; getAge({"name": "Alice", "age": 24}) + getAge({"name": "Alice", "age": 24});`,
			48,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
