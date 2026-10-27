package httpx

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

// OptionalParam is a query parameter that may be absent. Huma does not allow pointers in
// parameters, so the wrapper records whether the request set the value. Use Ptr to get a
// nil-able value for the layers below the handler.
//
// This follows the ParamWrapper and ParamReactor pattern from the Huma docs.
type OptionalParam[T any] struct {
	Value T
	IsSet bool
}

func (o OptionalParam[T]) Schema(r huma.Registry) *huma.Schema {
	return huma.SchemaFromType(r, reflect.TypeOf(o.Value))
}

func (o *OptionalParam[T]) Receiver() reflect.Value {
	return reflect.ValueOf(o).Elem().Field(0)
}

func (o *OptionalParam[T]) OnParamSet(isSet bool, _ any) {
	o.IsSet = isSet
}

// Ptr returns nil when the parameter was not set.
func (o OptionalParam[T]) Ptr() *T {
	if !o.IsSet {
		return nil
	}
	v := o.Value
	return &v
}
