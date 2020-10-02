package flatten

import "fmt"

type Merge struct {
	source []*Flatten
}

func NewMerge(flattens ...*Flatten) *Merge {
	return &Merge{
		source: flattens,
	}
}

func (m *Merge) Compile() *Flatten {
	if len(m.source) == 1 {
		return m.source[0]
	}

	sourceLen := len(m.source)
	result := *NewFlatten()

	for index, source := range m.source {
		if index == 0 && (index+1) >= sourceLen {
			result = mergeValue(*source, *m.source[index+1])
		} else if index == 1 {
			continue
		} else {
			result = mergeValue(result, *source)
		}
	}

	return &result
}

func mergeValue(left Flatten, right Flatten) Flatten {
	for rightKey, rightVal := range right.All(true) {
		left.Add(rightKey, rightVal)
	}

	fmt.Println(left.All(false))

	return left
}
