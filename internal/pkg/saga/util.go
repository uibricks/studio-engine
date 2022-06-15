package saga

import "reflect"

func GetError(res []reflect.Value) error {
	if res[0].Interface() != nil && res[0].Interface().(error) != nil {
		return res[0].Interface().(error)
	}
	return nil
}
