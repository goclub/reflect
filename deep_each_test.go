package xreflect

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)
func TestDeepEach1PnaicAndError(t *testing.T) {
	func(){
		defer func() {
			assert.Equal(t, recover(), errors.New("greject.DeepEach(&v, callback) callback can not be nil"))
		}()
		err := DeepEach1(struct{}{}, nil)
		assert.NoError(t, err)
	}()
}
func TestDeepEach1(t *testing.T) {
	type Info struct {
		FieldName string
		TypeName string
		TypeKind reflect.Kind
		AnonymousField bool
		JSONString string
		JSONTag string
		CanNotSet bool
	}
	var infos []Info
	type ID string
	type AnonymousCombination struct {
		Title string
	}
	type News struct {
		Content string
	}
	type IntList []int
	type Demo struct {
		Name string `json:"name"`
		Age uint
		UserID ID
		AnonymousCombination
		News News
		Hobby []string
		Array [1]string
		Numbers IntList
		NewsList []News
		StringPtr *string
		Map map[string]string
		MapPtr map[string]*string
		NewsPtr1 *News
		NewsPtr2 *News
		NewsListPtr *[]News
		NewsList2Ptr []*News
	}
	testStr := "orange"
	strElem := "ptr"
	strPtr := &strElem
	demo := Demo{
		Name:"nimoc", Age: 27,
		UserID: "a",
		AnonymousCombination: AnonymousCombination {
			Title: "t",
		},
		News: News{Content:"c"},
		Hobby: []string{"read"},
		Array: [1]string{"a"},
		Numbers: IntList{1},
		NewsList: []News{News{}},
		StringPtr: &testStr,
		Map: map[string]string{
			"type": "pass",
		},
		MapPtr: map[string]*string{
			"type": strPtr,
		},
		NewsPtr1: nil,
		NewsPtr2: &News{Content:""},
		NewsListPtr: &[]News{{Content:"a"}},
		NewsList2Ptr: []*News{{Content:"b"}},
	}
	actualInfos := []Info{
		{
			FieldName: "Name",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"nimoc"`,
			JSONTag: "name",
		},
		{
			FieldName: "Age",
			TypeName: "uint",
			TypeKind: reflect.Uint,
			JSONString: `27`,
		},
		{
			FieldName: "UserID",
			TypeName: "ID",
			TypeKind: reflect.String,
			JSONString: `"a"`,
		},
		{
			FieldName: "AnonymousCombination",
			TypeName: "AnonymousCombination",
			TypeKind: reflect.Struct,
			AnonymousField: true,
			JSONString: `{"Title":"t"}`,
		},
		{
			FieldName: "Title",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"t"`,
		},
		{
			FieldName: "News",
			TypeName: "News",
			TypeKind: reflect.Struct,
			AnonymousField: false,
			JSONString: `{"Content":"c"}`,
		},
		{
			FieldName: "Content",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"c"`,
		},
		{
			FieldName: "Hobby",
			TypeName: "",
			TypeKind: reflect.Slice,
			JSONString: `["read"]`,
		},
		{
			FieldName: "",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"read"`,
		},
		{
			FieldName: "Array",
			TypeName: "",
			TypeKind: reflect.Array,
			JSONString: `["a"]`,
		},
		{
			FieldName: "Numbers",
			TypeName: "IntList",
			TypeKind: reflect.Slice,
			JSONString: `[1]`,
		},
		{
			FieldName: "",
			TypeName: "int",
			TypeKind: reflect.Int,
			JSONString: `1`,
		},
		{
			FieldName: "NewsList",
			TypeName: "",
			TypeKind: reflect.Slice,
			JSONString: `[{"Content":""}]`,
		},
		{
			FieldName: "",
			TypeName: "News",
			TypeKind: reflect.Struct,
			JSONString: `{"Content":""}`,
		},
		{
			FieldName: "Content",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `""`,
		},
		{
			FieldName: "StringPtr",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"orange"`,
		},
		{
			FieldName: "Map",
			TypeName: "",
			TypeKind: reflect.Map,
			JSONString: `{"type":"pass"}`,
		},
		{
			FieldName: "",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"pass"`,
			CanNotSet: true,
		},
		{
			FieldName: "MapPtr",
			TypeName: "",
			TypeKind: reflect.Map,
			JSONString: `{"type":"ptr"}`,
		},
		{
			FieldName: "",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"ptr"`,
		},
		{
			FieldName: "NewsPtr2",
			TypeName: "News",
			TypeKind: reflect.Struct,
			JSONString: `{"Content":""}`,
		},
		{
			FieldName: "Content",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `""`,
		},
		{
			FieldName: "NewsListPtr",
			TypeName: "",
			TypeKind: reflect.Slice,
			JSONString: `[{"Content":"a"}]`,
		},
		{
			FieldName: "",
			TypeName: "News",
			TypeKind: reflect.Struct,
			JSONString: `{"Content":"a"}`,
		},
		{
			FieldName: "Content",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"a"`,
		},
		{
			FieldName: "NewsList2Ptr",
			TypeName: "",
			TypeKind: reflect.Slice,
			JSONString: `[{"Content":"b"}]`,
		},
		{
			FieldName: "",
			TypeName: "News",
			TypeKind: reflect.Struct,
			JSONString: `{"Content":"b"}`,
		},
		{
			FieldName: "Content",
			TypeName: "string",
			TypeKind: reflect.String,
			JSONString: `"b"`,
		},
	}
	err := DeepEach1(&demo, func(rValue reflect.Value, rType reflect.Type, field reflect.StructField) (op EachOperator) {
		infos = append(infos, Info{
			FieldName: field.Name,
			TypeName: rType.Name(),
			TypeKind: rType.Kind(),
			AnonymousField: field.Anonymous,
			JSONString: func() string {
				if rValue.CanInterface() {
					data, err := json.Marshal(rValue.Interface()) ; if err != nil {
						panic(err)
					}
					return string(data)
				} else {
					return "nil"
				}
			}(),
			JSONTag: field.Tag.Get("json"),
			CanNotSet: !rValue.CanSet(),
		})
		return op.Continue()
	})
	assert.NoError(t, err)
	assert.Equal(t, infos, actualInfos)
	if t.Failed() {
		log.Print(infos)
	}
}

func TestEachOperator(t *testing.T) {
	list := []string{"a","b","c"}
	msg := ""
	err := DeepEach1(&list, func(rValue reflect.Value, rType reflect.Type, field reflect.StructField) (op EachOperator) {
		msg += rValue.String()
		if rValue.String() == "b" {
			return op.Break()
		}
		return op.Continue()
	})
	assert.NoError(t, err)
	assert.Equal(t, msg, "ab")
}

func TestDeepEach1Map(t *testing.T) {
	type Item struct {
		Value string
	}

	type Data map[string]*Item

	data := map[string]*Item{"name":{"nimoc"}}
	err := DeepEach1(data, func(rValue reflect.Value, rType reflect.Type, field reflect.StructField) (op EachOperator) {
		if rType.Kind() == reflect.String {
			rValue.SetString(rValue.String() + "!")
		}
		return op.Continue()
	})
	assert.NoError(t ,err)
	assert.Equal(t, data, map[string]*Item{"name":{"nimoc!"}})
}
func TestDeepEach1Error(t *testing.T) {
	err := DeepEach1(map[string]int{"a":1,"b":2}, func(rValue reflect.Value, rType reflect.Type, field reflect.StructField) (op EachOperator) {
		if rValue.Int() == 2 {
			return op.Error(errors.New("value can not be 2"))
		}
		return op.Continue()
	})
	assert.Equal(t, err.Error(), "value can not be 2")
}

func TestDeepEach1Tree(t *testing.T) {
	type Node struct {
		Name string
		Value string
		Node *Node
	}
	node := Node{
		Name:"a",
		Value: "",
		Node: &Node{
			Name:"b",
			Value: "",
			Node: nil,
		},
	}
	var names []string
	err := DeepEach1(&node, func(rValue reflect.Value, rType reflect.Type, field reflect.StructField) (op EachOperator) {
		names  = append(names, field.Name)
		return op.Continue()
	})
	assert.NoError(t, err)
	assert.Equal(t, names, []string{"Name", "Value","Node", "Name", "Value"})
}