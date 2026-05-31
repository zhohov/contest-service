package domain

type TestCase struct {
	ID       int
	Input    string
	Expected string
}

func NewTestCase(input, excpected string) *TestCase {
	return &TestCase{
		Input:    input,
		Expected: excpected,
	}
}
