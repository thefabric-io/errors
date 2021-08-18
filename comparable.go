package errors

type Comparable interface {
	Compare(err1, err2 error) bool
}

func CompareMessageOnlyStrategy() Comparable {
	return compareMessageOnlyStrategy{}
}

type compareMessageOnlyStrategy struct{}

func (c compareMessageOnlyStrategy) Compare(err1, err2 error) bool {
	errA := Stack(err1)
	errB := Stack(err2)

	for _, a := range errA.stacks {
		for _, b := range errB.stacks {
			if b.message.Equal(a.message) {
				return true
			}
		}
	}

	return false
}

func CompareCodeOnlyStrategy() Comparable {
	return compareCodeOnlyStrategy{}
}

type compareCodeOnlyStrategy struct{}

func (c compareCodeOnlyStrategy) Compare(err1, err2 error) bool {
	errA := Stack(err1)
	errB := Stack(err2)

	for _, a := range errA.stacks {
		for _, b := range errB.stacks {
			if b.code.Equal(a.code) {
				return true
			}
		}
	}

	return false
}

func CompareStrictStrategy() Comparable {
	return compareStrictStrategy{}
}

type compareStrictStrategy struct{}

func (c compareStrictStrategy) Compare(err1, err2 error) bool {
	errA := Stack(err1)
	errB := Stack(err2)

	for _, a := range errA.stacks {
		for _, b := range errB.stacks {
			if b.code.Equal(a.code) && b.message.Equal(a.message) {
				return true
			}
		}
	}

	return false
}
