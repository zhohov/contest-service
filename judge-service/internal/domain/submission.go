package domain

type SubmissionStatus string

const (
	SubmissionStatePending SubmissionStatus = "pending"
	SubmissionStatsDone    SubmissionStatus = "done"
)

type SubmissionResult string

const (
	SubmissionResultUnknown      SubmissionResult = "unknown"
	SubmissionResultAccepted     SubmissionResult = "accepted"
	SubmissionResultWrongAnswer  SubmissionResult = "wrong_answer"
	SubmissionResultRuntimeError SubmissionResult = "runtime_error"
)

type Submission struct {
	ID     int
	Code   string
	Status SubmissionStatus
	Result SubmissionResult
}

// func NewSubmission(user User, problem Problem) *Submission {
// 	return &Submission{
// 		State:      SubmissionStatePending,
// 		StatusCode: SubmissionResultUnknown,
// 	}
// }

// func NewSubmission(code string, problem Problem) *Submission {
// 	return &Submission{
// 		Code:    code,
// 		Problem: problem,
// 		Status:  SubmissionStatePending,
// 		Result:  SubmissionResultUnknown,
// 	}
// }
