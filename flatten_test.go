package flatten

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var BenchmarkResult interface{}

func TestFlatten_Add(t *testing.T) {
	flat := NewFlatten()

	flat.Add("test.abs.1", "yep")

	assert.Equal(t, flat.Get("test.abs.1"), "yep")
}

func TestFlatten_Delete(t *testing.T) {
	flat := NewFlatten()

	flat.Add("test.abs.1", "yep")

	assert.Equal(t, flat.Get("test.abs.1"), "yep")

	flat.Delete("test.abs.1")

	assert.Equal(t, flat.Get("test.abs.1"), nil)
}

func TestFlatten_Get(t *testing.T) {
	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"a": map[string]interface{}{
				"avt": "hi",
			},
		},
	}

	data, err := NewFlattenFromMap(testedMap)

	if err != nil {
		t.Error(err)
	}

	assert.IsType(t, &Flatten{}, data.Get("test01.a"))
	assert.Equal(t, data.Get("test"), "a")
	assert.Equal(t, data.Get("test01.a.avt"), "hi")
}

func TestFlatten_ToJson(t *testing.T) {
	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"avt": "hi",
		},
		"abc": []interface{}{
			map[string]interface{}{
				"qw": 1,
			},
			map[string]interface{}{
				"qw": 2,
			},
			map[string]interface{}{
				"qw": 3,
			},
		},
		"dca": []interface{}{
			[]interface{}{
				1,2,3,
			},
			[]interface{}{
				4,5,6,
			},
		},
	}

	data, err := NewFlattenFromMap(testedMap)

	if err != nil {
		t.Error(err)
	}

	testedDataJson, err := json.Marshal(testedMap)

	if err != nil {
		t.Error(err)
	}


	assert.Equal(t, data.ToJson(), string(testedDataJson))
}

func TestFlatten_ToNested(t *testing.T) {
	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"avt": "hi",
		},
		"abc": []interface{}{
			map[string]interface{}{
				"qw": 1,
			},
			map[string]interface{}{
				"qw": 2,
			},
			map[string]interface{}{
				"qw": 3,
			},
		},
		"dca": []interface{}{
			[]interface{}{
				1,2,3,
			},
			[]interface{}{
				4,5,6,
			},
		},
	}

	data, err := NewFlattenFromMap(testedMap)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, data.ToNested(), testedMap)
}

func TestNewFlatten(t *testing.T) {
	assert.IsTypef(t, NewFlatten(), &Flatten{}, "")
}

func TestNewFlattenFromJson(t *testing.T) {
	testedString := "{\"test\": \"a\",\"test01\": {\"a\": {\"avt\": \"hi\"}}\n}"

	data, err := NewFlattenFromJson(testedString)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, data.Get("test"), "a")
	assert.IsType(t, &Flatten{}, data.Get("test01.a"))
}

func TestNewFlattenFromMapFlatten(t *testing.T) {
	testedMap := map[string]interface{}{
		"test":         "a",
		"test01.a.avt": "hi",
	}

	data, err := NewFlattenFromMap(testedMap)

	if err != nil {
		t.Error(err)
	}

	assert.IsType(t, &Flatten{}, data.Get("test01.a"))
	assert.Equal(t, data.Get("test"), "a")
	assert.Equal(t, data.Get("test01.a.avt"), "hi")
}

func TestNewFlattenFromMapNested(t *testing.T) {
	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"a": map[string]interface{}{
				"avt": "hi",
			},
		},
	}

	data, err := NewFlattenFromMap(testedMap)

	if err != nil {
		t.Error(err)
	}

	assert.IsType(t, &Flatten{}, data.Get("test01.a"))
	assert.Equal(t, data.Get("test"), "a")
	assert.Equal(t, data.Get("test01.a.avt"), "hi")
}

func BenchmarkNewFlatten(b *testing.B) {
	var r *Flatten

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r = NewFlatten()
	}

	BenchmarkResult = r
}

func BenchmarkNewFlattenFromJson(b *testing.B) {
	var r *Flatten
	testData := "{\n\t\"test\": \"a\",\n\t\"test01\": {\n\t\t\"a\": {\n\t\t\t\"avt\": \"hi\"\n\t\t}\n\t},\n\t\"abc\": [{\n\t\t\t\"qw\": 1\n\t\t},\n\t\t{\n\t\t\t\"qw\": 2\n\t\t},\n\t\t{\n\t\t\t\"qw\": 3\n\t\t}\n\t],\n\t\"dca\": [\n\t\t[1, 2, 3],\n\t\t[4, 5, 6]\n\t]\n}"

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r, _ = NewFlattenFromJson(testData)
	}

	BenchmarkResult = r
}

func BenchmarkNewFlattenFromMap(b *testing.B) {
	var r *Flatten

	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"avt": "hi",
		},
		"abc": []interface{}{
			map[string]interface{}{
				"qw": 1,
			},
			map[string]interface{}{
				"qw": 2,
			},
			map[string]interface{}{
				"qw": 3,
			},
		},
		"dca": []interface{}{
			[]interface{}{
				1,2,3,
			},
			[]interface{}{
				4,5,6,
			},
		},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r, _ = NewFlattenFromMap(testedMap)
	}

	BenchmarkResult = r
}

func BenchmarkFlatten_Get(b *testing.B) {
	var r interface{}

	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"avt": "hi",
		},
		"abc": []interface{}{
			map[string]interface{}{
				"qw": 1,
			},
			map[string]interface{}{
				"qw": 2,
			},
			map[string]interface{}{
				"qw": 3,
			},
		},
		"dca": []interface{}{
			[]interface{}{
				1,2,3,
			},
			[]interface{}{
				4,5,6,
			},
		},
	}

	flatten, _ := NewFlattenFromMap(testedMap)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r = flatten.Get("abc.0.qw")
	}

	BenchmarkResult = r
}

func BenchmarkFlatten_Get_ParentPart(b *testing.B) {
	var r interface{}

	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"avt": "hi",
		},
		"abc": []interface{}{
			map[string]interface{}{
				"qw": 1,
			},
			map[string]interface{}{
				"qw": 2,
			},
			map[string]interface{}{
				"qw": 3,
			},
		},
		"dca": []interface{}{
			[]interface{}{
				1,2,3,
			},
			[]interface{}{
				4,5,6,
			},
		},
	}

	flatten, _ := NewFlattenFromMap(testedMap)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r = flatten.Get("abc.0")
	}

	BenchmarkResult = r
}

func BenchmarkFlatten_ToNested(b *testing.B) {
	var r interface{}

	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"avt": "hi",
		},
		"abc": []interface{}{
			map[string]interface{}{
				"qw": 1,
			},
			map[string]interface{}{
				"qw": 2,
			},
			map[string]interface{}{
				"qw": 3,
			},
		},
		"dca": []interface{}{
			[]interface{}{
				1,2,3,
			},
			[]interface{}{
				4,5,6,
			},
		},
	}

	flatten, _ := NewFlattenFromMap(testedMap)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r = flatten.ToNested()
	}

	BenchmarkResult = r
}

func BenchmarkFlatten_ToJson(b *testing.B) {
	var r interface{}

	testedMap := map[string]interface{}{
		"test": "a",
		"test01": map[string]interface{}{
			"avt": "hi",
		},
		"abc": []interface{}{
			map[string]interface{}{
				"qw": 1,
			},
			map[string]interface{}{
				"qw": 2,
			},
			map[string]interface{}{
				"qw": 3,
			},
		},
		"dca": []interface{}{
			[]interface{}{
				1,2,3,
			},
			[]interface{}{
				4,5,6,
			},
		},
	}

	flatten, _ := NewFlattenFromMap(testedMap)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r = flatten.ToJson()
	}

	BenchmarkResult = r
}

func BenchmarkNewFlattenFromJson_BigData(b *testing.B) {
	var r *Flatten

	testingJson, err := ioutil.ReadFile("./test_1.json")

	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r, _ = NewFlattenFromJson(string(testingJson))
	}

	BenchmarkResult = r
}

func BenchmarkNewFlattenFromMap_BigData(b *testing.B) {
	var r *Flatten

	testingJson, err := ioutil.ReadFile("./test_1.json")

	if err != nil {
		b.Error(err)
	}

	testingData := new(interface{})

	err = json.Unmarshal(testingJson, testingData)

	if err != nil {
		b.Error(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r, _ = NewFlattenFromMap(*testingData)
	}

	BenchmarkResult = r
}

func BenchmarkFlatten_Get_BigData(b *testing.B) {
	var r interface{}

	testingJson, err := ioutil.ReadFile("./test_1.json")

	if err != nil {
		b.Error(err)
	}

	testingData := new(interface{})

	err = json.Unmarshal(testingJson, testingData)

	if err != nil {
		b.Error(err)
	}

	flatten, _ := NewFlattenFromMap(*testingData)

	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r = flatten.Get("text.login_screen_footer_text")
	}

	BenchmarkResult = r
}

func BenchmarkFlatten_Get_ParentBigData(b *testing.B) {

	var r interface{}

	testingJson, err := ioutil.ReadFile("./test_1.json")

	if err != nil {
		b.Error(err)
	}

	testingData := new(interface{})

	err = json.Unmarshal(testingJson, testingData)

	if err != nil {
		b.Error(err)
	}

	flatten, _ := NewFlattenFromMap(*testingData)


	b.ReportAllocs()
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		r = flatten.Get("text")
	}

	BenchmarkResult = r
}
