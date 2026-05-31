package domain

type SubmissionStatus string

const (
	SubmissionStatusPending SubmissionStatus = "pending"
	SubmissionStatusDone    SubmissionStatus = "done"
)

type SubmissionResult string

const (
	SubmissionResultUnknown      SubmissionResult = "unknown"
	SubmissionResultAccepted     SubmissionResult = "accepted"
	SubmissionResultWrongAnswer  SubmissionResult = "wrong_answer"
	SubmissionResultRuntimeError SubmissionResult = "runtime_error"
)

type Submission struct {
	ID      int
	User    User
	Code    string
	Problem Problem
	Status  SubmissionStatus
	Result  SubmissionResult
}

// func NewSubmission(user User, problem Problem) *Submission {
// 	return &Submission{
// 		State:      SubmissionStatePending,
// 		StatusCode: SubmissionResultUnknown,
// 	}
// }

func NewSubmission(code string, user User, problem Problem) *Submission {
	return &Submission{
		User:    user,
		Code:    code,
		Problem: problem,
		Status:  SubmissionStatusPending,
		Result:  SubmissionResultUnknown,
	}
}
