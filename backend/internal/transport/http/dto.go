package http

type CreateArticleRequestJson struct {
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
}

type CreateArticleResponseJson struct {
	Article ArticleJson `json:"article"`
}

type ArticleJson struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	Content string `json:"content"`
}

type GetArticlesJson struct {
	Articles []ArticleJson `json:"articles"`
}

type TestCaseJson struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type ProblemJson struct {
	ID         int            `json:"id"`
	Title      string         `json:"title"`
	Slug       string         `json:"slug"`
	Difficulty string         `json:"difficulty"`
	Content    string         `json:"content"`
	TestCases  []TestCaseJson `json:"test_cases"`
}

type CreateProblemRequestJson struct {
	Title      string         `json:"title"`
	Slug       string         `json:"slug"`
	Difficulty string         `json:"difficulty"`
	Content    string         `json:"content"`
	TestCases  []TestCaseJson `json:"test_cases"`
}

type SubmissionJson struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Status string `json:"status"`
	Result string `json:"result"`
}

type CreateProblemResponseJson struct {
	Problem ProblemJson `json:"problem"`
}

type GetProblemsJson struct {
	Problems []ProblemJson `json:"problems"`
}

type LoginRequestJson struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponseJson struct {
	Token string `json:"token"`
}

type RegisterRequestJson struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponseJson struct {
	Token string `json:"token"`
}

type CreateSubmissionRequestJson struct {
	ProblemID int    `json:"problem_id"`
	Code      string `json:"code"`
}

type CreateSubmissionResponseJson struct {
	SubmissionID int `json:"submission_id"`
}

// type GetSubmissionByIdResponseJson struct {
// 	Submission SubmissionJson
// }

type CreateContestRequestJson struct {
	Title       string                     `json:"title"`
	Slug        string                     `json:"slug"`
	Description string                     `json:"description"`
	StartsAt    string                     `json:"starts_at"`
	EndsAt      string                     `json:"ends_at"`
	Problems    []CreateProblemRequestJson `json:"problems"`
}

type CreateContestResponseJson struct {
	Contest ContestJson `json:"contest"`
}

type ContestJson struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Slug        string        `json:"slug"`
	Description string        `json:"description"`
	CreatorID   int           `json:"creator_id"`
	StartsAt    string        `json:"starts_at"`
	EndsAt      string        `json:"ends_at"`
	Problems    []ProblemJson `json:"problems"`
}

type GetContestBySlugResponseJson struct {
	Contest ContestJson `json:"contest"`
}
