package fp

import (
	"errors"
	"reflect"
)

// Map ...
// "function"
// INPUTs: (val T), (val T, idx int), (val T, idx int, list []T)
// OUTPUTs: (T'), (T', error)
func Map(slice, function interface{}) MapResult {
	vslice := reflect.ValueOf(slice)

	switch vslice.Kind() {
	case reflect.Array, reflect.Slice:
	default:
		panic("input is not an array or slice")
	}

	tfunction := reflect.TypeOf(function)

	if tfunction.Kind() != reflect.Func {
		panic("function must be a func")
	}
	if tfunction.NumIn() == 0 || tfunction.NumIn() > 3 {
		panic("function must have at least 1 and at most 3 arguments")
	}
	if tfunction.NumOut() != 1 && tfunction.NumOut() != 2 {
		panic("function must have at least 1 and at most 2 outputs")
	}

	ret := reflect.MakeSlice(reflect.SliceOf(tfunction.Out(0)), 0, vslice.Len())

	if vslice.Len() == 0 {
		return MapResult{val: ret}
	}

	idx := 0

	input := make([]reflect.Value, tfunction.NumIn())
	for i := 0; i < tfunction.NumIn(); i++ {
		switch i {
		case 0:
			input[i] = vslice.Index(idx)
		case 1:
			input[i] = reflect.ValueOf(idx)
		case 2:
			input[i] = vslice
		}
	}

	for i := idx; i < vslice.Len(); i++ {
		input[0] = vslice.Index(i)
		if len(input) > 1 {
			input[1] = reflect.ValueOf(i)
		}

		out := reflect.ValueOf(function).Call(input)

		switch len(out) {
		case 1:
			ret = reflect.Append(ret, out[0])
		case 2:
			err, ok := out[1].Interface().(error)
			if !ok {
				return MapResult{err: errors.New("fp: function second output must be of error type")}
			}

			return MapResult{err: err}
		default:
			return MapResult{err: errors.New("fp: function call return a number of outputs different from 1 or 2")}
		}
	}

	return MapResult{val: ret.Interface()}
}

type MapResult struct {
	val interface{}
	err error
}

func (r MapResult) Get() interface{} {
	return r.val
}

func (r MapResult) Error() error {
	return r.err
}

func (r MapResult) Map(function interface{}) MapResult {
	if r.err != nil {
		return r
	}

	return Map(r.val, function)
}
