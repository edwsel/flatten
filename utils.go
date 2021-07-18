package flatten

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

func flatten(flatMap map[string]interface{}, nested interface{}, key string, delimiter string) error {
	var rNested reflect.Value

	assign := func(newKey string, value reflect.Value) error {
		switch value.Kind() {
		case reflect.Map, reflect.Array, reflect.Slice:
			if err := flatten(flatMap, value.Interface(), newKey, delimiter); err != nil {
				return err
			}
		case reflect.Interface:
			switch value.Interface().(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, string, []byte:
				flatMap[newKey] = value.Interface()
			default:
				if err := flatten(flatMap, value.Interface(), newKey, delimiter); err != nil {
					return err
				}
			}
		default:
			flatMap[newKey] = value.Interface()
		}

		return nil
	}

	if val, ok := nested.(reflect.Value); ok {
		rNested = val
	} else {
		rNested = reflect.ValueOf(nested)
	}

	switch rNested.Kind() {
	case reflect.Map:
		rangeRNested := rNested.MapRange()
		for rangeRNested.Next() {
			k := rangeRNested.Key()
			v := rangeRNested.Value()

			if k.Kind() != reflect.String {
				return errors.New("key type must be a string")
			}

			if v.Kind() == reflect.Interface {
				v = reflect.ValueOf(v.Interface())
			}

			newKey := makeFlattenKey(key, k.String(), delimiter)
			err := assign(newKey, v)

			if err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < rNested.Len(); i++ {
			v := rNested.Index(i)

			newKey := makeFlattenKey(key, strconv.Itoa(i), delimiter)

			if v.Kind() == reflect.Interface {
				v = reflect.ValueOf(v.Interface())
			}

			err := assign(newKey, v)

			if err != nil {
				return err
			}
		}
	default:
		return NotValidInputError
	}

	return nil
}

func nested(flat map[string]interface{}, delimiter string) interface{} {
	var result interface{}

	for key, value := range flat {
		subKeys := strings.Split(key, delimiter)
		result = makeNested(result, subKeys, value)
	}

	return result
}

func makeNested(nested interface{}, subKeys []string, value interface{}) interface{} {
	var isLast = len(subKeys) == 1
	var isArray bool
	var indexArray int

	if num, err := strconv.Atoi(subKeys[0]); err == nil {
		isArray = true
		indexArray = num
	}

	if nested == nil && isArray {
		nested = make([]interface{}, indexArray+1)
	} else if nested == nil && !isArray {
		nested = make(map[string]interface{}, 0)
	} else if typedNested, ok := nested.([]interface{}); ok && (indexArray+1) > len(typedNested) {
		tmpNestedArray := make([]interface{}, indexArray+1)

		copy(tmpNestedArray, typedNested)

		nested = tmpNestedArray
	}

	switch nested.(type) {
	case []interface{}:
		if isLast {
			nested.([]interface{})[indexArray] = value
		} else {
			nested.([]interface{})[indexArray] = makeNested(nested.([]interface{})[indexArray], subKeys[1:], value)
		}
		break
	case map[string]interface{}:
		if isLast {
			nested.(map[string]interface{})[subKeys[0]] = value
		} else {
			nested.(map[string]interface{})[subKeys[0]] = makeNested(nested.(map[string]interface{})[subKeys[0]], subKeys[1:], value)
		}
		break
	}

	return nested
}

func makeFlattenKey(key string, subKey string, delimiter string) string {
	if key == "" {
		return subKey
	}

	return key + delimiter + subKey
}

func makeKey(prefix string, key string) string {
	if len(prefix) > 0 {
		key = prefix + "." + key
	}

	return key
}
