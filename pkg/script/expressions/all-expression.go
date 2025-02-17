package expressions

import (
	"fmt"
	"github.com/threagile/threagile/pkg/script/common"
)

type AllExpression struct {
	literal    string
	in         common.ValueExpression
	item       string
	index      string
	expression common.BoolExpression
}

func (what *AllExpression) ParseBool(script any) (common.BoolExpression, any, error) {
	what.literal = common.ToLiteral(script)

	switch script.(type) {
	case map[string]any:
		for key, value := range script.(map[string]any) {
			switch key {
			case common.In:
				item, errorExpression, itemError := new(ValueExpression).ParseValue(value)
				if itemError != nil {
					return nil, errorExpression, fmt.Errorf("failed to parse %q of any-expression: %v", key, itemError)
				}

				what.in = item

			case common.Item:
				text, ok := value.(string)
				if !ok {
					return nil, value, fmt.Errorf("failed to parse %q of any-expression: expected string, got %T", key, value)
				}

				what.item = text

			case common.Index:
				text, ok := value.(string)
				if !ok {
					return nil, value, fmt.Errorf("failed to parse %q of any-expression: expected string, got %T", key, value)
				}

				what.index = text

			default:
				if what.expression != nil {
					return nil, script, fmt.Errorf("failed to parse any-expression: additional bool expression %q", key)
				}

				expression, errorScript, itemError := new(ExpressionList).ParseAny(map[string]any{key: value})
				if itemError != nil {
					return nil, errorScript, fmt.Errorf("failed to parse any-expression: %v", itemError)
				}

				boolExpression, ok := expression.(common.BoolExpression)
				if !ok {
					return nil, script, fmt.Errorf("any-expression contains non-bool expression: %v", itemError)
				}

				what.expression = boolExpression
			}
		}

	default:
		return nil, script, fmt.Errorf("failed to parse any-expression: expected map[string]any, got %T", script)
	}

	return what, nil, nil
}

func (what *AllExpression) ParseAny(script any) (common.Expression, any, error) {
	return what.ParseBool(script)
}

func (what *AllExpression) EvalBool(scope *common.Scope) (bool, string, error) {
	oldIterator := scope.SwapIterator(nil)
	defer scope.SetIterator(oldIterator)

	inValue, errorEvalLiteral, evalError := what.in.EvalAny(scope)
	if evalError != nil {
		return false, errorEvalLiteral, evalError
	}

	switch castValue := inValue.(type) {
	case []any:
		if what.expression == nil {
			return true, "", nil
		}

		for index, item := range castValue {
			if len(what.index) > 0 {
				scope.Set(what.index, index)
			}

			scope.SetIterator(item)
			if len(what.item) > 0 {
				scope.Set(what.item, item)
			}

			value, errorLiteral, expressionError := what.expression.EvalBool(scope)
			if expressionError != nil {
				return false, errorLiteral, fmt.Errorf("error evaluating expression #%v of any-expression: %v", index+1, expressionError)
			}

			if !value {
				return false, "", nil
			}
		}

	case map[string]any:
		if what.expression == nil {
			return true, "", nil
		}

		for name, item := range castValue {
			if len(what.index) > 0 {
				scope.Set(what.index, name)
			}

			scope.SetIterator(item)
			if len(what.item) > 0 {
				scope.Set(what.item, item)
			}

			value, errorLiteral, expressionError := what.expression.EvalBool(scope)
			if expressionError != nil {
				return false, errorLiteral, fmt.Errorf("error evaluating expression %q of any-expression: %v", name, expressionError)
			}

			if !value {
				return false, "", nil
			}
		}

	default:
		return false, what.Literal(), fmt.Errorf("failed to eval any-expression: expected iterable type, got %T", inValue)
	}

	return true, "", nil
}

func (what *AllExpression) EvalAny(scope *common.Scope) (any, string, error) {
	return what.EvalBool(scope)
}

func (what *AllExpression) Literal() string {
	return what.literal
}
