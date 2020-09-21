package flatten

import (
	"strconv"
	"strings"
)

func flatten(flatMap map[string]interface{}, nested interface{}, key string) error {
	assign := func(newKey string, value interface{}) error {
		switch value.(type) {
		case map[string]interface{}:
			if err := flatten(flatMap, value.(map[string]interface{}), newKey); err != nil {
				return err
			}
		case []interface{}:
			if err := flatten(flatMap, value.([]interface{}), newKey); err != nil {
				return err
			}
		default:
			flatMap[newKey] = value
		}

		return nil
	}

	switch nested.(type) {
	case map[string]interface{}:
		for k, v := range nested.(map[string]interface{}) {
			newKey := makeFlattenKey(key, k)
			err := assign(newKey, v)

			if err != nil {
				return err
			}
		}
	case []interface{}:
		for k, v := range nested.([]interface{}) {
			newKey := makeFlattenKey(key, strconv.Itoa(k))

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

func nested(flat map[string]interface{}) interface{} {
	var result interface{}

	for key, value := range flat {
		subKeys := strings.Split(key, ".")
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

func makeNested2(nested interface{}, subKeys []string, value interface{}) interface{} {
	for _, subKey := range subKeys {

	}


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

func makeFlattenKey(key string, subKey string) string {
	if key == "" {
		return subKey
	}

	return key + "." + subKey
}
