package evaluator

import (
	"github.com/hlongvu/monkeylang/monkey/object"
	"testing"
)

func TestEvalBooleanExpression(t *testing.T){
	tests := []struct{
		input string
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

		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests{
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool{
	result, ok := obj.(*object.Boolean)
	if !ok{
		t.Errorf("Object is not boolean, got %T", obj)
		return false
	}

	if result.Value != expected{
		t.Errorf("Object has wrong value, got %t, want %t", result.Value, expected)
		return false
	}
	return true
}