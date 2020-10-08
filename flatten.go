package flatten

import (
	"encoding/json"
	"errors"
	"strings"
)

const (
	DEFAULT_DELIMITER = "."
)

var (
	NotValidInputError = errors.New("not valid input error")
)

type Flatten struct {
	namespace string
	delimiter string
	container map[string]interface{}
	keyStore  map[string][]string
}

func NewFlatten() *Flatten {
	return &Flatten{
		namespace: "",
		delimiter: DEFAULT_DELIMITER,
		container: map[string]interface{}{},
		keyStore:  map[string][]string{},
	}
}

func NewFlattenFromMap(data interface{}, delimiter string) (*Flatten, error) {
	result := NewFlatten()
	result.delimiter = delimiter

	err := flatten(result.container, data, "", delimiter)

	if err != nil {
		return nil, err
	}

	for key := range result.container {
		result.metaKeyAdd(key)
	}

	return result, nil
}

func NewFlattenFromJson(data string, delimiter string) (*Flatten, error) {
	dataMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(data), &dataMap)

	if err != nil {
		return nil, err
	}

	result := NewFlatten()

	err = flatten(result.container, dataMap, "", delimiter)

	if err != nil {
		return nil, err
	}

	for key := range result.container {
		result.metaKeyAdd(key)
	}

	return result, nil
}

func (f *Flatten) GetDelimiter() string {
	return f.delimiter
}

func (f *Flatten) SetDelimiter(delimiter string) *Flatten {
	f.delimiter = delimiter

	return f
}

func (f *Flatten) GetNamespace() string {
	return f.namespace
}

func (f *Flatten) SetNamespace(prefix string) *Flatten {
	f.namespace = prefix

	return f
}

func (f *Flatten) Add(key string, value interface{}) *Flatten {
	sliceKey := strings.Split(key, f.delimiter)

	for index, _ := range sliceKey {
		subKey := strings.Join(sliceKey[:index+1], f.delimiter)
		if _, ok := f.keyStore[subKey]; ok && len(f.keyStore[subKey]) == 1 {
			f.Delete(subKey)
		}
	}

	switch value.(type) {
	case map[interface{}]interface{}, []interface{}:
		newData := map[string]interface{}{}
		err := flatten(newData, value, key, f.delimiter)

		if err != nil {
			break
		}

		for newKey, newVal := range newData {
			f.container[newKey] = newVal
		}
	default:
		f.container[key] = value
	}

	f.metaKeyAdd(key)

	return f
}

func (f *Flatten) Has(key string) bool {
	return false
}

func (f *Flatten) Get(key string) interface{} {
	flat := NewFlatten()

	if _, ok := f.keyStore[key]; !ok {
		return NewFlatten()
	}

	if val, ok := f.container[key]; len(f.keyStore[key]) == 1 && ok {
		return val
	}

	for _, value := range f.keyStore[key] {
		flat.Add(value[len(key)+1:], f.container[value])
	}

	return flat
}

func (f *Flatten) All(withNamespace bool) map[string]interface{} {
	result := map[string]interface{}{}

	for key, value := range f.container {
		if len(f.namespace) > 0 && withNamespace {
			key = f.namespace + f.delimiter + key
		}

		result[key] = value
	}

	return result
}

func (f Flatten) Delete(key string) {
	key = makeKey(f.namespace, key)

	delete(f.container, key)
	f.metaKeyDelete(key)

}

func (f *Flatten) ToJson(withNamespace bool) string {
	data, err := json.Marshal(f.ToNested(withNamespace))

	if err != nil {
		return ""
	}

	return string(data)
}

func (f *Flatten) ToNested(withNamespace bool) interface{} {
	return nested(f.All(withNamespace), f.delimiter)
}

func (f *Flatten) metaKeyAdd(key string) {
	subKeys := strings.Split(key, f.delimiter)

	if len(subKeys) == 1 && !f.metaKeyExist(subKeys[0], key) {
		f.keyStore[subKeys[0]] = []string{key}
	}

	for index := range subKeys {
		sliceSubKey := strings.Join(subKeys[:index+1], f.delimiter)

		if !f.metaKeyExist(sliceSubKey, key) {
			f.keyStore[sliceSubKey] = append(f.keyStore[sliceSubKey], key)
		}
	}
}

func (f *Flatten) metaKeyExist(subKey string, key string) bool {
	if _, ok := f.keyStore[subKey]; ok {
		for _, val := range f.keyStore[subKey] {
			if val == key {
				return true
			}
		}
	}

	return false
}

func (f *Flatten) metaKeyDelete(key string) {
	subKeys := strings.Split(key, f.delimiter)

	if len(subKeys) == 1 && !f.metaKeyExist(subKeys[0], key) {
		f.keyStore[subKeys[0]] = []string{key}
	}

	for index := range subKeys {
		sliceSubKey := strings.Join(subKeys[:index+1], f.delimiter)

		if _, ok := f.keyStore[sliceSubKey]; ok {
			for index, value := range f.keyStore[sliceSubKey] {
				if value != key {
					continue
				}

				f.keyStore[sliceSubKey] = append(f.keyStore[sliceSubKey][:index], f.keyStore[sliceSubKey][index+1:]...)

				break
			}

			if len(f.keyStore[sliceSubKey]) == 0 {
				delete(f.keyStore, sliceSubKey)
			}
		}
	}
}
