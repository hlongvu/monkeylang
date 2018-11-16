package evaluator

import (
	"fmt"
	"github.com/hlongvu/monkeylang/monkey/ast"
	"github.com/hlongvu/monkeylang/monkey/object"
)

var (
	NULL = &object.Null{}
	TRUE = &object.Boolean{Value:true}
	FALSE = &object.Boolean{Value:false}
)

func Eval(node ast.Node, env *object.Environment) object.Object{
	switch node := node.(type) {

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.Boolean:
		return nativeToBooleanObject(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right){
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left){
			return left
		}
		right := Eval(node.Right, env)
		if isError(right){
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val){
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val){
			return val
		}

		env.Set(node.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{
			Parameters: params,
			Body:body,
			Env:env,
		}

	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function){
			return function
		}

		args := evalExpression(node.Arguments, env)
		if len(args) == 1 && isError(args[0]){
			return args[0]
		}

		return applyFunction(function, args)

	}
	return nil
}

func applyFunction(fn object.Object, args []object.Object) object.Object{
 	function, ok := fn.(*object.Function)
 	if !ok{
 		return newError("Not a function: %s", fn.Type())
	}
	extendedEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment{
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramId, param := range fn.Parameters{
		env.Set(param.Value, args[paramId])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object{
	if returnValue, ok := obj.(*object.ReturnValue); ok{
		return returnValue.Value
	}
	return obj
}

func evalExpression(exps []ast.Expression, env *object.Environment) []object.Object{
	var result []object.Object
	for _, e := range exps{
		evaluated := Eval(e, env)
		if isError(evaluated){
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}


func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object{
	val, ok := env.Get(node.Value)
	if !ok{
		return newError("identifier not found: "+ node.Value)
	}else{
		return val
	}
}

func newError(format string, a ...interface{}) *object.Error{
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object)bool{
	if obj!=nil{
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object  {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition){
		return Eval(ie.Consequence, env)
	}else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}else{
		return NULL
	}
}

func isTruthy(obj object.Object) bool{
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}


func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object{
	switch {
	case left.Type() == object.INTEGER_OBJ  && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return nativeToBooleanObject(left == right)
	case operator == "!=":
		return nativeToBooleanObject(left != right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",  left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object{
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}

	case "<":
		return nativeToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeToBooleanObject(leftVal != rightVal)

	default:
		return newError("unknown operator: %s %s %s",  left.Type(), operator, right.Type())
	}
}


func evalPrefixExpression(operator string, right object.Object) object.Object{
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
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
func evalMinusPrefixOperatorExpression(right object.Object) object.Object{
	if right.Type() != object.INTEGER_OBJ{
		return newError("unknown operator: -%s", right.Type())
	}

	value:=right.(*object.Integer).Value
	return &object.Integer{Value: - value}
}


func evalProgram( sts [] ast.Statement,  env *object.Environment) object.Object{
	var result object.Object
	for _, st := range sts{
		result = Eval(st, env)
		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatements(block *ast.BlockStatement,  env *object.Environment) object.Object{
	var result object.Object
	for _, stmt := range block.Statements{
		result  = Eval(stmt, env)
		if result != nil && (result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ) {
			return result
		}
	}
	return result
}

func nativeToBooleanObject(input bool) *object.Boolean{
	if input{
		return TRUE
	}else{
		return FALSE
	}
}

