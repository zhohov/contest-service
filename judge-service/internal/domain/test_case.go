package domain

type TestCase struct {
	Input    string
	Expected string
}

func NewTestCase(input, excpected string) *TestCase {
	return &TestCase{
		Input:    input,
		Expected: excpected,
	}
}
