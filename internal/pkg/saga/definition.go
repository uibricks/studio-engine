package saga

import (
	"golang.org/x/net/context"
	"reflect"

	"fmt"
)

func fn(ctx context.Context) {
}

type subTxDefinitions map[string]subTxDefinition

type subTxDefinition struct {
	subTxID    string
	action     reflect.Value
	compensate reflect.Value
}

func (s subTxDefinitions) addDefinition(subTxID string, action interface{}, compensate interface{}) subTxDefinitions {
	actionMethod := subTxMethod(action)
	compensateMethod := reflect.ValueOf(nil)
	if compensate != nil {
		compensateMethod = subTxMethod(compensate)
	}
	s[subTxID] = subTxDefinition{
		subTxID:    subTxID,
		action:     actionMethod,
		compensate: compensateMethod,
	}
	return s
}

func (s subTxDefinitions) findDefinition(subTxID string) (subTxDefinition, bool) {
	define, ok := s[subTxID]
	return define, ok
}

type paramTypeRegister struct {
	nameToType map[string]reflect.Type
	typeToName map[reflect.Type]string
}

func (r *paramTypeRegister) addParams(fc interface{}) {
	funcValue := subTxMethod(fc)

	funcType := funcValue.Type()
	for i := 0; i < funcType.NumIn(); i++ {
		paramType := funcType.In(i)
		pTypName := paramType.Name()
		if pTypName == "" {
			pTypName = fmt.Sprintf("%s", paramType)
		}
		r.nameToType[pTypName] = paramType
		r.typeToName[paramType] = pTypName
	}
	for i := 0; i < funcType.NumOut(); i++ {
		returnType := funcType.Out(i)
		r.nameToType[returnType.Name()] = returnType
		r.typeToName[returnType] = returnType.Name()
	}

}

func (r *paramTypeRegister) findTypeName(typ reflect.Type) (string, bool) {

	ctx := context.Background()
	ctx = context.WithValue(ctx, "a", "b")
	paramType := reflect.TypeOf(ctx)

	if paramType == typ {
		funcValue := reflect.ValueOf(fn)
		funcType := funcValue.Type()
		paramType := funcType.In(0)
		typ = paramType
	}

	f, ok := r.typeToName[typ]
	return f, ok
}

func (r *paramTypeRegister) findType(typeName string) (reflect.Type, bool) {
	f, ok := r.nameToType[typeName]
	return f, ok
}

func subTxMethod(obj interface{}) reflect.Value {

	funcValue := reflect.ValueOf(obj)

	if funcValue.Kind() != reflect.Func {
		panic("Register object must be a func")
	}
	if funcValue.Type().NumIn() < 1 ||
		funcValue.Type().In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
		panic("First argument must use context.Context.")
	}
	return funcValue
}
