package evaluator

import (
	"github.com/hlongvu/monkeylang/monkey/object"
	"testing"
)

func TestErrorhandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
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
			`
			if (10 > 1) {
  				if (10 > 1) {
   					 return true + false;
 				 }

  				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN"},
	}

	for _, tt := range tests{
		evaluated := testEval(tt.input)
		errorObj, ok := evaluated.(*object.Error)
		if !ok{
			t.Errorf("No error return, got %T", evaluated)
			continue
		}
		if errorObj.Message != tt.expectedMessage{
			t.Errorf("Wrong error message, expected %q, got %q", tt.expectedMessage, errorObj.Message)
			continue
		}
	}
}
