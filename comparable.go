package errors

type Comparable interface {
	Compare(err1, err2 error) bool
}

func oneIsNil(err1, err2 error) bool {
	return err1 == nil || err2 == nil
}

func CompareMessageOnlyStrategy() Comparable {
	return compareStrategy{
		tester: func(err1, err2 Error) bool {
			return err2.message.Equal(err1.message)
		},
	}
}

func CompareCodeOnlyStrategy() Comparable {
	return compareStrategy{
		tester: func(err1, err2 Error) bool {
			return err2.code.Equal(err1.code)
		},
	}
}

func CompareStrictStrategy() Comparable {
	return compareStrategy{
		tester: func(err1, err2 Error) bool {
			return err2.code.Equal(err1.code) && err2.message.Equal(err1.message)
		},
	}
}

type equalityTester func(err1, err2 Error) bool

type compareStrategy struct {
	tester equalityTester
}

func (c compareStrategy) Compare(err1, err2 error) bool {
	if oneIsNil(err1, err2) {
		return err1 == err2
	}

	errA := Stack(err1).(*Errors)
	errB := Stack(err2).(*Errors)

	for _, a := range errA.stacks {
		for _, b := range errB.stacks {
			if c.tester(a, b) {
				return true
			}
		}
	}

	return false
}
