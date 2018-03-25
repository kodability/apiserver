package run

// TryoutMockRunner is a mock class for TryoutRunner interface
type TryoutMockRunner struct {
	Result *JUnitReport
	Err    error
}

func (r *TryoutMockRunner) Run(lang, code, testCode string) (*JUnitReport, error) {
	return r.Result, r.Err
}
