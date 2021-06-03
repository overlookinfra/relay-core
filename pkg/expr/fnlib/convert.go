package fnlib

import (
	"context"
	"reflect"

	"github.com/puppetlabs/relay-core/pkg/expr/convert"
	"github.com/puppetlabs/relay-core/pkg/expr/fn"
	"github.com/puppetlabs/relay-core/pkg/expr/model"
)

var convertMarkdownDescriptor = fn.DescriptorFuncs{
	DescriptionFunc: func() string { return "Converts a string in markdown format to another applicable syntax" },
	PositionalInvokerFunc: func(ev model.Evaluator, args []interface{}) (fn.Invoker, error) {
		if len(args) != 2 {
			return nil, &fn.ArityError{Wanted: []int{2}, Variadic: false, Got: len(args)}
		}

		fn := fn.EvaluatedPositionalInvoker(ev, args, func(ctx context.Context, args []interface{}) (m interface{}, err error) {
			to, found := args[0].(string)
			if !found {
				return nil, &fn.PositionalArgError{
					Arg: 0,
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf("")},
						Got:    reflect.TypeOf(args[0]),
					},
				}
			}

			switch md := args[1].(type) {
			case string:
				r, err := convert.ConvertMarkdown(convert.ConvertType(to), []byte(md))
				if err != nil {
					return nil, err
				}
				return string(r), nil
			default:
				return nil, &fn.PositionalArgError{
					Arg: 1,
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf("")},
						Got:    reflect.TypeOf(args[1]),
					},
				}
			}
		})
		return fn, nil
	},
	KeywordInvokerFunc: func(ev model.Evaluator, args map[string]interface{}) (fn.Invoker, error) {
		for _, arg := range []string{"to", "content"} {
			if _, found := args[arg]; !found {
				return nil, &fn.KeywordArgError{Arg: arg, Cause: fn.ErrArgNotFound}
			}
		}

		return fn.EvaluatedKeywordInvoker(ev, args, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
			ct, ok := args["to"].(string)
			if !ok {
				return nil, &fn.KeywordArgError{
					Arg: "to",
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf("")},
						Got:    reflect.TypeOf(args["to"]),
					},
				}
			}

			switch md := args["content"].(type) {
			case string:
				r, err := convert.ConvertMarkdown(convert.ConvertType(ct), []byte(md))
				if err != nil {
					return nil, err
				}
				return string(r), nil
			default:
				return nil, &fn.KeywordArgError{
					Arg: "content",
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf("")},
						Got:    reflect.TypeOf(md),
					},
				}
			}
		}), nil
	},
}
