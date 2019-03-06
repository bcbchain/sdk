package mapx

import (
	"fmt"
	"reflect"
	"sort"
)

// ForRange - range map object
func ForRange(mapObj interface{}, f interface{}) {

	// check map object
	mapObjValue := reflect.ValueOf(mapObj)
	mapObjType := reflect.TypeOf(mapObj)
	keyType := mapObjType.Key()
	valueType := mapObjType.Elem()

	if mapObjType.Kind() != reflect.Map {
		panic("mapObj not type of map")
	}

	// check operation function
	typeOfF := reflect.TypeOf(f)
	numIn := typeOfF.NumIn()
	if numIn != 2 {
		panic("f must be 2 in parameters")
	}

	if typeOfF.In(0) != keyType {
		panic(fmt.Sprintf("f's first in parameter's type should be %s, obtain %s",
			keyType.String(), typeOfF.In(0).String()))
	}

	if typeOfF.In(1) != valueType {
		panic(fmt.Sprintf("f's second in parameter's type should be %s, obtain %s",
			valueType.String(), typeOfF.In(1).String()))
	}

	// sort keys
	ks := mapObjValue.MapKeys()
	sort.SliceStable(ks, func(i, j int) bool {
		switch keyType.Kind() {
		case reflect.String:
			return ks[i].String() < ks[j].String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return ks[i].Int() < ks[j].Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return ks[i].Uint() < ks[j].Uint()
		case reflect.Float32, reflect.Float64:
			return ks[i].Float() < ks[j].Float()
		case reflect.Bool:
			return ks[i].Bool() == false && ks[j].Bool() == true
		default:
			panic(fmt.Sprintf("do not support key's type:%s", keyType.String()))
		}

		return false
	})

	// range map object
	for _, k := range ks {
		fValue := reflect.ValueOf(f)

		in := make([]reflect.Value, 2)
		in[0] = k
		in[1] = mapObjValue.MapIndex(k)

		fValue.Call(in)
	}
}
