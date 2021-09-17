package errors

type Comparable interface {
	Compare(err1, err2 error) bool
}

func oneIsNil(err1, err2 error) bool {
	return err1 == nil || err2 == nil
}

func CompareMessageOnlyStrategy() Comparable {
	return compareMessageOnlyStrategy{}
}

type compareMessageOnlyStrategy struct{}

func (c compareMessageOnlyStrategy) Compare(err1, err2 error) bool {
	if oneIsNil(err1, err2) {
		return err1 == err2
	}

	errA := Stack(err1).(*Errors)
	errB := Stack(err2).(*Errors)

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
	if oneIsNil(err1, err2) {
		return err1 == err2
	}

	errA := Stack(err1).(*Errors)
	errB := Stack(err2).(*Errors)

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
	if oneIsNil(err1, err2) {
		return err1 == err2
	}
	errA := Stack(err1).(*Errors)
	errB := Stack(err2).(*Errors)

	for _, a := range errA.stacks {
		for _, b := range errB.stacks {
			if b.code.Equal(a.code) && b.message.Equal(a.message) {
				return true
			}
		}
	}

	return false
}
