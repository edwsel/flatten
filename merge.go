package flatten

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

	result := NewFlatten()

	for _, source := range m.source {
		result = mergeValue(result, source)
	}

	return result
}

func mergeValue(left *Flatten, right *Flatten) *Flatten {
	for rightKey, rightVal := range right.All(true) {
		left.Add(rightKey, rightVal)
	}

	return left
}
