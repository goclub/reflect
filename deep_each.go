package xreflect

import (
	"errors"
	"reflect"
)

type DeepEachInfo struct {
	IsRoot bool
}
// 前序遍历，遇到指针跳过指针回调值
type EachCallback func(rValue reflect.Value, rType reflect.Type, field reflect.StructField) (op EachOperator)
type EachOperator struct {
	shouldBreak bool
	error error
}
func (op EachOperator) Error(err error) EachOperator {
	if err == nil {
		panic(errors.New("greflect: warning Error(err) err is nil, maybe you write wrong code."))
	}
	return EachOperator{
		error: err,
		shouldBreak: true,
	}
}
func (op EachOperator) Continue() EachOperator {
	return EachOperator{}
}
func (op EachOperator) Break() EachOperator {
	return EachOperator{
		shouldBreak: true,
	}
}
func (op EachOperator) shouldReturn() bool {
	if op.error != nil {
		return true
	}
	if op.shouldBreak {
		return true
	}
	return false
}
func deepEach(v interface{}, callback EachCallback, write bool) error {
	if callback == nil {
		panic(errors.New("greject.DeepEach(&v, callback) callback can not be nil"))
	}

	rootValue := reflect.ValueOf(v)
	rootType := rootValue.Type()

	if write {
		if !rootValue.CanSet() && rootType.Kind() != reflect.Map {
			if rootValue.Kind() == reflect.Ptr {
				rootValue = rootValue.Elem()
				rootType = rootType.Elem()
			}
			if !rootValue.CanSet() {
				return errors.New("DeepEach(v, callback) v must can set, mu be you should use DeepEach(&v, callback)")
			}
		}
	}
	info := DeepEachInfo{
		IsRoot: true,
	}
	op := coreEach(coreEachProps{
		parentValue: rootValue,
		parentType: rootType,
		field: reflect.StructField{},
		callback: callback,
		info: info,
	})
	return op.error
}
func OnlyReadDeepEach1(v interface{}, callback EachCallback) error {
	return deepEach(v, callback, false)
}
func DeepEach1(v interface{}, callback EachCallback) error {
	return deepEach(v, callback, true)
}
type coreEachProps struct {
	parentValue reflect.Value
	parentType reflect.Type
	field reflect.StructField
	callback EachCallback
	info DeepEachInfo
}
func coreEach(props coreEachProps) EachOperator {
	switch {
	case props.info.IsRoot:
	case props.parentType.Kind() == reflect.Ptr:
	default:
		op := props.callback(props.parentValue, props.parentType, props.field)
		if op.shouldReturn() { return op }
	}
	if props.info.IsRoot {
		props.info.IsRoot = false
	}
	switch props.parentType.Kind() {
	case reflect.Ptr:
		if !props.parentValue.IsNil() {
			elementValue := props.parentValue.Elem()
			op := coreEach(coreEachProps{
				parentValue: elementValue,
				parentType:  elementValue.Type(),
				field:       props.field,
				callback:    props.callback,
				info:        props.info,
			})
			if op.shouldReturn() { return op }
		}
	case reflect.Map:
		for _, key := range props.parentValue.MapKeys() {
			mapValue := props.parentValue.MapIndex(key)
			op := coreEach(coreEachProps{
				parentValue: mapValue,
				parentType:  mapValue.Type(),
				field:       reflect.StructField{},
				callback:    props.callback,
				info:        props.info,
			})
			if op.shouldReturn() { return op }
		}
	case reflect.Struct:
		for i:=0;i< props.parentType.NumField();i++ {
			rValue := props.parentValue.Field(i)
			rType := rValue.Type()
			field := props.parentType.Field(i)
			op := coreEach(coreEachProps{
				parentValue: rValue,
				parentType: rType,
				field: field,
				callback: props.callback,
				info: props.info,
			})
			if op.shouldReturn() { return op }
		}
	case reflect.Slice:
		for i:=0;i<props.parentValue.Len();i++ {
			elementValue := props.parentValue.Index(i)
			elementType := elementValue.Type()
			op := coreEach(coreEachProps{
				parentValue: elementValue,
				parentType: elementType,
				field: reflect.StructField{},
				callback: props.callback,
				info: props.info,
			})
			if op.shouldReturn() { return op }
		}
	default:
		// ignore other type
	}
	return EachOperator{}
}