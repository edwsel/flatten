package flatten

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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
		delimiter: ".",
		container: map[string]interface{}{},
		keyStore:  map[string][]string{},
	}
}

func NewFlattenFromMap(data interface{}) (*Flatten, error) {
	result := NewFlatten()

	err := flatten(result.container, data, "")

	if err != nil {
		return nil, err
	}

	for key := range result.container {
		result.metaKeyAdd(key)
	}

	return result, nil
}

func NewFlattenFromJson(data string) (*Flatten, error) {
	dataMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(data), &dataMap)

	if err != nil {
		return nil, err
	}

	result := NewFlatten()

	err = flatten(result.container, dataMap, "")

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
	fmt.Println(f.container, f.container[key], key, value)

	f.container[key] = value

	f.metaKeyAdd(key)

	return f
}

func (f *Flatten) Get(key string) interface{} {
	if _, ok := f.keyStore[key]; !ok {
		return nil
	}

	if val, ok := f.container[key]; len(f.keyStore[key]) == 1 && ok {
		return val
	}

	flat := NewFlatten()

	for _, value := range f.keyStore[key] {
		flat.Add(value[len(key)+1:], f.container[value])
	}

	return flat
}

func (f *Flatten) All(withNamespace bool) map[string]interface{} {
	result := map[string]interface{}{}

	for key, value := range f.container {
		if len(f.namespace) > 0 && withNamespace {
			key = f.namespace + "." + key
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
	return nested(f.All(withNamespace))
}

func (f *Flatten) metaKeyAdd(key string) {
	subKeys := strings.Split(key, ".")

	if len(subKeys) == 1 && !f.metaKeyExist(subKeys[0], key) {
		f.keyStore[subKeys[0]] = []string{key}
	}

	for index := range subKeys {
		sliceSubKey := strings.Join(subKeys[:index+1], ".")

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
	subKeys := strings.Split(key, ".")

	if len(subKeys) == 1 && !f.metaKeyExist(subKeys[0], key) {
		f.keyStore[subKeys[0]] = []string{key}
	}

	for index := range subKeys {
		sliceSubKey := strings.Join(subKeys[:index+1], ".")

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
