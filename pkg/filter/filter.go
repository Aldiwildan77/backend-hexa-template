package filter

import "reflect"

func RemoveNulls(m map[string]interface{}) map[string]interface{} {
	val := reflect.ValueOf(m)
	for _, e := range val.MapKeys() {
		v := val.MapIndex(e)
		if v.IsNil() {
			delete(m, e.String())
			continue
		}

		switch t := v.Interface().(type) {
		case map[string]interface{}:
			RemoveNulls(t)
		}
	}

	return m
}
