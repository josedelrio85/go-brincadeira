package functions

import (
	"fmt"
	"reflect"
	"strconv"
)

//For a defined Struct, with json annnotations
func GetJsonKeysReceived(source interface{}) ([]string, []string) {
	s := reflect.ValueOf(source).Elem()
	k := reflect.TypeOf(source).Elem()
	campos := []string{}
	values := []string{}
	typeOfT := s.Type()

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		// fmt.Printf("%d: %s %s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface())

		//if somo of the properties of the struct is empty, obtain its equivalent as tag json key
		//and add to array. The finally is to obtain an array with the received keys
		// if f.Interface() != "" {
		zz, ok := k.FieldByName(typeOfT.Field(i).Name)
		if ok {
			campos = append(campos, string(zz.Tag.Get("json")))

			var val string
			switch f.Kind() {
			case reflect.Float64:
				val = fmt.Sprintf("%.0f", f.Interface().(float64))
			case reflect.Int:
				val = strconv.FormatInt(f.Interface().(int64), 'i')
			case reflect.String:
				val = f.Interface().(string)
			}
			values = append(values, val)
			// values = append(values, f.Interface().(string))
			// values = append(values, "a")
		}
		// }
	}
	return campos, values
}
