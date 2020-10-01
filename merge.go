package flatten

type Merge struct {
	source []*Flatten
}

func NewMerge(flattens ...*Flatten) *Merge {
	return &Merge{
		source: flattens,
	}
}

func (m *Merge) Compile() Flatten {
	if len(m.source) == 1 {
		return *m.source[0]
	}

	sourceLen := len(m.source)
	result := Flatten{}

	for key, source := range m.source {
		if key == 0 && (key+1) >= sourceLen {
			result = mergeValue(*source, *m.source[key+1])
		} else if key == 1 {
			continue
		} else {
			result = mergeValue(result, *source)
		}
	}

	return result
}

func mergeValue(left Flatten, right Flatten) Flatten {
	for key, rightVal := range right.All(true) {
		left.Add(key, rightVal)
	}

	return left
}
