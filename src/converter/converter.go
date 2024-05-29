package converter

import (
	"log"
	"reflect"
	"strings"
)

func TransformStructToMapJson(val reflect.Value, noEmpty bool) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	tpe := val.Type()
	for i := 0; i < tpe.NumField(); i++ {
		f := tpe.Field(i)
		v, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}
		if val.FieldByName(f.Name).Kind() == reflect.Slice {
			arr, _ := TransformSliceToMapJson(val.FieldByName(f.Name), noEmpty)
			res[strings.Split(v, ",")[0]] = arr

		} else if val.FieldByName(f.Name).Kind() == reflect.Struct {
			arr, _ := TransformStructToMapJson(val.FieldByName(f.Name), noEmpty)
			res[strings.Split(v, ",")[0]] = arr

		} else if val.FieldByName(f.Name).Kind() == reflect.Pointer {
			// val.FieldByName(f.Name)
			// reflect.New(val.FieldByName(f.Name).Elem()).Elem()
			newEl := val.FieldByName(f.Name)
			if newEl.IsNil() {
				newEl = reflect.New(val.FieldByName(f.Name).Type().Elem())
			}
			arr, _ := TransformStructToMapJson(newEl.Elem(), noEmpty)
			res[strings.Split(v, ",")[0]] = arr

		} else {
			res[strings.Split(v, ",")[0]] = val.FieldByName(f.Name).Interface()
		}
	}
	return res, nil
}

func TransformSliceToMapJson(val reflect.Value, noEmpty bool) ([]interface{}, error) {
	res := []interface{}{}
	if val.Len() == 0 {
		if noEmpty {
			if val.Type().Elem().Kind() == reflect.Pointer {
				val = reflect.MakeSlice(reflect.SliceOf(val.Type().Elem().Elem()), 1, 1)
			} else {
				val = reflect.MakeSlice(reflect.SliceOf(val.Type().Elem()), 1, 1)
			}
		} else {
			return []interface{}{}, nil
		}
	}
	// v := val.Index(0)
	for i := 0; i < val.Len(); i++ {
		if val.Index(i).Kind() == reflect.Struct {
			it, _ := TransformStructToMapJson(val.Index(i), noEmpty)
			res = append(res, it)

		} else if val.Index(i).Kind() == reflect.Slice {
			it, _ := TransformSliceToMapJson(val.Index(i), noEmpty)
			res = append(res, it)
		} else if val.Index(i).Kind() == reflect.Pointer {
			log.Println(val.Index(i), val.Type().Elem())
			reflect.New(val.Type().Elem().Elem()).Elem()
			it, _ := TransformStructToMapJson(val.Index(i).Elem(), noEmpty)
			log.Println(it)
			res = append(res, it)

		} else {
			res = append(res, val.Index(i).Interface())
		}
	}
	return res, nil
}

func TransformMapJsonToStruct(myMap map[string]interface{}, val reflect.Value) (reflect.Value, error) {
	tpe := val.Type()
	res := reflect.New(tpe).Elem()
	for i := 0; i < tpe.NumField(); i++ {
		f := tpe.Field(i)
		v, ok := f.Tag.Lookup("json")
		if !ok {
			continue
		}
		jsonName := strings.Split(v, ",")[0]
		vl, ok := myMap[jsonName]
		if !ok || vl == nil {
			continue
		}
		if reflect.ValueOf(vl).Kind() == reflect.Slice {
			intf := []interface{}{}
			for j := 0; j < reflect.ValueOf(vl).Len(); j++ {
				intf = append(intf, reflect.ValueOf(vl).Index(j).Interface())
			}
			rs, _ := TransformSliceToStruct(intf, res.FieldByName(f.Name).Type().Elem(), reflect.ValueOf(vl).Len())
			res.FieldByName(f.Name).Set(rs.Convert(f.Type))
		} else if reflect.ValueOf(vl).Kind() == reflect.Map {
			myNextMap := map[string]interface{}{}
			for _, e := range reflect.ValueOf(vl).MapKeys() {
				va := reflect.ValueOf(vl).MapIndex(e)
				key := e.String()
				myNextMap[key] = va.Interface()
			}

			newEl := reflect.New(val.FieldByName(f.Name).Type().Elem())

			rs, _ := TransformMapJsonToStruct(myNextMap, newEl.Elem())
			res.FieldByName(f.Name).Set(rs.Addr().Convert(f.Type))
		} else {
			res.FieldByName(f.Name).Set(reflect.ValueOf(vl).Convert(f.Type))
		}
	}
	return res, nil

}

func TransformSliceToStruct(arr []interface{}, tpe reflect.Type, length int) (reflect.Value, error) {
	res := reflect.Value{}
	if tpe.Kind() == reflect.Pointer {
		res = reflect.MakeSlice(reflect.SliceOf(tpe.Elem()), length, length)
	} else {
		res = reflect.MakeSlice(reflect.SliceOf(tpe), length, length)
	}

	for i, it := range arr {
		vl := res.Index(i)
		if vl.Kind() == reflect.Map {
			myNextMap := map[string]interface{}{}
			for _, e := range vl.MapKeys() {
				va := vl.MapIndex(e)
				key := e.String()
				myNextMap[key] = va.Interface()
			}

			newEl := reflect.New(vl.Type().Elem())
			rs, _ := TransformMapJsonToStruct(myNextMap, newEl.Elem())
			vl.Set(rs.Addr().Convert(vl.Type()))
		} else {
			vl.Set(reflect.ValueOf(it).Convert(tpe))
		}
	}
	return res, nil

}
