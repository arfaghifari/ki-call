package converter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var myMapExample = map[string]interface{}{}

// Transform Json Map to Struct
func TransformMapJsonToStruct(myMapInterface interface{}, valueType reflect.Type) (interface{}, error) {
	var err error
	var temp interface{}
	var fieldValue, value reflect.Value
	if reflect.TypeOf(myMapInterface) != reflect.TypeOf(myMapExample) {
		return nil, fmt.Errorf("value error: not a valid json format. want: %v, got %v", reflect.TypeOf(myMapExample), reflect.TypeOf(myMapInterface))
	}
	myMapReflect := reflect.ValueOf(myMapInterface)
	myMap := map[string]interface{}{}
	for _, key := range myMapReflect.MapKeys() {
		mapKey := key.String()
		myMap[mapKey] = myMapReflect.MapIndex(key).Interface()
	}
	if valueType.Kind() != reflect.Struct && valueType.Kind() != reflect.Pointer {
		return nil, errors.New("value error: not a struct")
	}
	if valueType.Kind() == reflect.Pointer {
		temp, err = TransformMapJsonToStruct(myMap, valueType.Elem())
		tempAddr := reflect.New(valueType.Elem())
		tempAddr.Elem().Set(reflect.ValueOf(temp))
		return tempAddr.Interface(), err
	}
	value = reflect.New(valueType).Elem()
	// should be struct
	for i := 0; i < valueType.NumField(); i++ {
		f := valueType.Field(i)
		v, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}
		jsonName := strings.Split(v, ",")[0]
		val, ok := myMap[jsonName]
		if !ok || val == nil {
			continue
		}
		fieldValue, err = TransformInterfaceToFieldValue(val, f.Type)
		if err != nil {
			return nil, err
		}
		value.FieldByName(f.Name).Set(fieldValue)
	}
	return value.Interface(), nil
}

func TransformInterfaceToFieldValue(val interface{}, valueType reflect.Type) (reflect.Value, error) {
	var err error
	var temp reflect.Value
	var temp2 interface{}
	var nilValue = reflect.ValueOf(nil)
	value := reflect.ValueOf(val)
	if value.Interface() == nil {
		return nilValue, nil
	}
	if valueType.Kind() == reflect.Pointer {
		temp, err = TransformInterfaceToFieldValue(val, valueType.Elem())
		tempAddr := reflect.New(valueType.Elem())
		tempAddr.Elem().Set(temp)
		return tempAddr, err
	}
	if valueType.Kind() == reflect.Struct {
		temp2, err = TransformMapJsonToStruct(val, valueType)
		return reflect.ValueOf(temp2), err
	}
	if valueType.Kind() == reflect.Slice {
		if value.Type().Kind() != reflect.Slice {
			return nilValue, errors.New("value error: not a slice")
		}
		res := reflect.MakeSlice(valueType, value.Len(), value.Len())
		var it reflect.Value
		for i := 0; i < value.Len(); i++ {
			it, err = TransformInterfaceToFieldValue(value.Index(i).Interface(), valueType.Elem())
			if err != nil {
				break
			}
			res.Index(i).Set(it)
		}
		return res, nil
	}
	if err != nil {
		return nilValue, err
	}

	return reflect.ValueOf(val).Convert(valueType), nil
}

// Transform Struct To Json Map
func TransformStructToMapJson(val interface{}, noEmpty bool) (map[string]interface{}, error) {
	var err error
	res := map[string]interface{}{}
	value := reflect.ValueOf(val)
	valueType := value.Type()
	if valueType.Kind() != reflect.Struct && valueType.Kind() != reflect.Pointer {
		return nil, errors.New("value error: not a struct")
	}
	if valueType.Kind() == reflect.Pointer {
		res, err = TransformStructToMapJson(reflect.Indirect(value).Interface(), noEmpty)
		return res, err
	}
	// should be struct
	for i := 0; i < valueType.NumField(); i++ {
		f := valueType.Field(i)
		v, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}
		res[strings.Split(v, ",")[0]], err = TransformInterfaceToMapJson(value.FieldByName(f.Name).Interface(), noEmpty)
		if err != nil {
			break
		}
	}
	return res, err
}

func TransformInterfaceToMapJson(val interface{}, noEmpty bool) (interface{}, error) {
	var err error
	value := reflect.ValueOf(val)
	valueType := value.Type()

	if valueType.Kind() == reflect.Pointer {
		if value.IsNil() {
			if noEmpty {
				newVal := reflect.New(valueType.Elem())
				return TransformInterfaceToMapJson(newVal.Elem().Interface(), noEmpty)
			}
			return nil, nil
		}
		return TransformInterfaceToMapJson(reflect.Indirect(value).Interface(), noEmpty)
	}

	if valueType.Kind() == reflect.Struct {
		return TransformStructToMapJson(val, noEmpty)
	}

	if valueType.Kind() == reflect.Slice {
		if value.IsNil() && !noEmpty {
			return nil, nil
		}
		if value.IsNil() || value.Len() == 0 {
			if noEmpty {
				var it interface{}
				newVal := reflect.New(valueType.Elem())
				it, err = TransformInterfaceToMapJson(newVal.Elem().Interface(), noEmpty)
				return []interface{}{it}, err
			}
		}
		res := []interface{}{}
		var it interface{}
		for i := 0; i < value.Len(); i++ {
			it, err = TransformInterfaceToMapJson(value.Index(i).Interface(), noEmpty)
			if err != nil {
				break
			}
			res = append(res, it)
		}
		return res, nil
	}
	if err != nil {
		return nil, err
	}

	return val, nil
}
