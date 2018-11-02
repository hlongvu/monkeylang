package evaluator

import (
	"github.com/hlongvu/monkeylang/monkey/object"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct{
		input string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests{
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool{
	result, ok := obj.(*object.Integer)
	if !ok{
		t.Errorf("Object is not Integer, got %T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected{
		t.Errorf("Object has wrong value. Got %d, want %d", result.Value, expected)
		return false
	}

	return true
}