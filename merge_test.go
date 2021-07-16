package flatten

import (
	"fmt"
	"testing"
)

func TestMerge_Compile(t *testing.T) {
	firstElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.0.test": "aaa",
	}, ".")

	secondElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.0.test": map[string]interface{}{
			"abc": "100",
		},
	}, ".")

	result := NewMerge(firstElem, secondElem).Compile()

	fmt.Printf("%+v\n", result.ToNested(true))
}

func Benchmark_Merge_Compile_2Map(b *testing.B) {
	firstElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "aaa",
		"root.2.test": "bbb",
		"root.3.test": "bbb",
	}, ".")

	secondElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "111",
		"root.4.test": "bbb",
		"root.3.test": "lll",
	}, ".")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewMerge(firstElem, secondElem).Compile()
	}
}

func Benchmark_Merge_Compile_3Map(b *testing.B) {
	//b.SetBytes(48)
	firstElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "aaa",
		"root.2.test": "bbb",
		"root.3.test": "bbb",
	}, ".")

	secondElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "111",
		"root.4.test": "bbb",
		"root.3.test": "lll",
	}, ".")

	thirdElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "111",
		"root.4.test": "bbb",
		"root.3.test": "lll",
	}, ".")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewMerge(firstElem, secondElem, thirdElem).Compile()
	}
}

func Benchmark_Merge_Compile_4Map(b *testing.B) {
	//b.SetBytes(48)
	firstElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "aaa",
		"root.2.test": "bbb",
		"root.3.test": "bbb",
	}, ".")

	secondElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "111",
		"root.4.test": "bbb",
		"root.3.test": "lll",
	}, ".")

	thirdElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "111",
		"root.4.test": "bbb",
		"root.3.test": "lll",
	}, ".")

	fourthElem, _ := NewFlattenFromMap(map[string]interface{}{
		"root.1.test": "111",
		"root.4.test": "bbb",
		"root.3.test": "lll",
	}, ".")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewMerge(firstElem, secondElem, thirdElem, fourthElem).Compile()
	}
}
