package http

import (
	"time"

	"github.com/zhohov/contest-service/backend/internal/usecase"
)

func mapCreateArticleToInput(req CreateArticleRequestJson) usecase.CreateArticleInput {
	return usecase.CreateArticleInput{
		Title:   req.Title,
		Slug:    req.Slug,
		Content: req.Content,
	}
}

func mapGetArticlesOutputToJson(out *usecase.GetArticlesOutput) GetArticlesJson {
	articles := make([]ArticleJson, 0, len(out.Articles))

	for _, article := range out.Articles {
		articles = append(articles, mapArticleDTOToJson(article))
	}

	return GetArticlesJson{
		Articles: articles,
	}
}

func mapCreateArticleOutputToJson(out *usecase.CreateArticleOutput) CreateArticleResponseJson {
	return CreateArticleResponseJson{
		Article: mapArticleDTOToJson(out.Article),
	}
}

func mapArticleDTOToJson(article usecase.ArticleDTO) ArticleJson {
	return ArticleJson{
		ID:      article.ID,
		Title:   article.Title,
		Slug:    article.Slug,
		Content: article.Content,
	}
}

func mapCreateProblemToInput(req CreateProblemRequestJson) usecase.CreateProblemInput {
	testCases := make([]usecase.TestCaseDTO, 0, len(req.TestCases))
	for _, testCase := range req.TestCases {
		testCases = append(testCases, usecase.TestCaseDTO{
			Input:    testCase.Input,
			Expected: testCase.Expected,
		})
	}

	return usecase.CreateProblemInput{
		Title:      req.Title,
		Slug:       req.Slug,
		Difficulty: req.Difficulty,
		Content:    req.Content,
		TestCases:  testCases,
	}
}

func mapProblemDTOToJson(problem usecase.ProblemDTO) ProblemJson {
	testCases := make([]TestCaseJson, 0, len(problem.TestCases))
	for _, testCase := range problem.TestCases {
		testCases = append(testCases, TestCaseJson{
			Input:    testCase.Input,
			Expected: testCase.Expected,
		})
	}

	return ProblemJson{
		ID:         problem.ID,
		Title:      problem.Title,
		Slug:       problem.Slug,
		Difficulty: problem.Difficulty,
		Content:    problem.Content,
		TestCases:  testCases,
	}
}

func mapCreateProblemOutputToJson(out *usecase.CreateProblemOutput) CreateProblemResponseJson {
	return CreateProblemResponseJson{
		Problem: mapProblemDTOToJson(out.Problem),
	}
}

func mapGetProblemsOutputToJson(out *usecase.GetProblemsOutput) GetProblemsJson {
	problems := make([]ProblemJson, 0, len(out.Problems))
	for _, problem := range out.Problems {
		problems = append(problems, mapProblemDTOToJson(problem))
	}

	return GetProblemsJson{
		Problems: problems,
	}
}

func mapGetProblemBySlugToInput(slug string) usecase.GetProblemBySlugInput {
	return usecase.GetProblemBySlugInput{
		Slug: slug,
	}
}

func mapGetProblemBySlugOutputToJson(out *usecase.GetProblemBySlugOutput) ProblemJson {
	return mapProblemDTOToJson(out.Problem)
}

func mapLoginToinput(req LoginRequestJson) usecase.LoginInput {
	return usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

func mapRegisterToInput(req RegisterRequestJson) usecase.RegisterInput {
	return usecase.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
}

func mapLoginToJson(out *usecase.LoginOutput) LoginResponseJson {
	return LoginResponseJson{
		Token: out.Token,
	}
}

func mapRegisterToJson(out *usecase.RegisterOutput) RegisterResponseJson {
	return RegisterResponseJson{
		Token: out.Token,
	}
}

func mapCreateSubmissionToInput(req CreateSubmissionRequestJson) usecase.CreateSubmissionInput {
	return usecase.CreateSubmissionInput{
		ProblemID: req.ProblemID,
		Code:      req.Code,
	}
}

func mapCreateSubmissionOutputToJson(out *usecase.CreateSubmissionOutput) CreateSubmissionResponseJson {
	return CreateSubmissionResponseJson{
		SubmissionID: out.SubmissionID,
	}
}

func mapGetSubmissionByIdToInput(id int) usecase.GetSubmissionByIdInput {
	return usecase.GetSubmissionByIdInput{
		ID: id,
	}
}

// func mapGetSubmissionByIdOutputToJson(out *usecase.GetSubmissionByIdOutput) GetSubmissionByIdResponseJson {
// 	return GetSubmissionByIdResponseJson{
// 		Submission: mapSubmissionDTOToJson(out.Submission),
// 	}
// }

func mapGetSubmissionByIdOutputToJson(out *usecase.GetSubmissionByIdOutput) SubmissionJson {
	return mapSubmissionDTOToJson(out.Submission)
}

func mapSubmissionDTOToJson(submission usecase.SubmissionDTO) SubmissionJson {
	return SubmissionJson{
		ID:     submission.ID,
		Code:   submission.Code,
		Status: submission.Status,
		Result: submission.Result,
	}
}

func mapCreateContestToInput(
	req CreateContestRequestJson,
) (usecase.CreateContestInput, error) {
	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		return usecase.CreateContestInput{}, err
	}

	endsAt, err := time.Parse(time.RFC3339, req.EndsAt)
	if err != nil {
		return usecase.CreateContestInput{}, err
	}

	problems := make([]usecase.CreateContestProblemInput, 0, len(req.Problems))

	for _, problem := range req.Problems {
		testCases := make([]usecase.CreateContestTestCaseInput, 0, len(problem.TestCases))

		for _, testCase := range problem.TestCases {
			testCases = append(testCases, usecase.CreateContestTestCaseInput{
				Input:    testCase.Input,
				Expected: testCase.Expected,
			})
		}

		problems = append(problems, usecase.CreateContestProblemInput{
			Title:      problem.Title,
			Slug:       problem.Slug,
			Difficulty: problem.Difficulty,
			Content:    problem.Content,
			TestCases:  testCases,
		})
	}

	return usecase.CreateContestInput{
		Title:       req.Title,
		Slug:        req.Slug,
		Description: req.Description,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		Problems:    problems,
	}, nil
}

func mapGetContestBySlugToInput(slug string) usecase.GetContestBySlugInput {
	return usecase.GetContestBySlugInput{
		Slug: slug,
	}
}

func mapGetContestBySlugOutputToJson(out *usecase.GetContestBySlugOutput) GetContestBySlugResponseJson {
	return GetContestBySlugResponseJson{
		Contest: mapContestDTOToJson(out.Contest),
	}
}

func mapCreateContestOutputToJson(out *usecase.CreateContestOutput) CreateContestResponseJson {
	return CreateContestResponseJson{
		Contest: mapContestDTOToJson(out.Contest),
	}
}

func mapContestDTOToJson(contest usecase.ContestDTO) ContestJson {
	problems := make([]ProblemJson, 0, len(contest.Problems))

	for _, problem := range contest.Problems {
		problems = append(problems, mapProblemDTOToJson(problem))
	}

	return ContestJson{
		ID:          contest.ID,
		Title:       contest.Title,
		Slug:        contest.Slug,
		Description: contest.Description,
		CreatorID:   contest.CreatorID,
		StartsAt:    contest.StartsAt.Format(time.RFC3339),
		EndsAt:      contest.EndsAt.Format(time.RFC3339),
		Problems:    problems,
	}
}
