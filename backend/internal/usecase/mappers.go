package usecase

import "github.com/zhohov/contest-service/backend/internal/domain"

func mapArticleToDTO(article *domain.Article) ArticleDTO {
	return ArticleDTO{
		ID:      article.ID,
		Title:   article.Title,
		Slug:    article.Slug,
		Content: article.Content,
	}
}

func mapArticlesToDTOs(articles []domain.Article) []ArticleDTO {
	articleDTOs := make([]ArticleDTO, 0, len(articles))
	for _, article := range articles {
		articleDTOs = append(articleDTOs, mapArticleToDTO(&article))
	}
	return articleDTOs
}

func mapTestCaseToDTO(testCase *domain.TestCase) TestCaseDTO {
	return TestCaseDTO{
		Input:    testCase.Input,
		Expected: testCase.Expected,
	}
}

func mapProblemToDTO(problem *domain.Problem) ProblemDTO {
	testCases := make([]TestCaseDTO, 0, len(problem.TestCases))
	for _, testCase := range problem.TestCases {
		testCases = append(testCases, mapTestCaseToDTO(&testCase))
	}

	return ProblemDTO{
		ID:         problem.ID,
		Title:      problem.Title,
		Slug:       problem.Slug,
		Difficulty: problem.Difficulty,
		Content:    problem.Content,
		TestCases:  testCases,
	}
}

func mapProblemsToDTOs(problems []domain.Problem) []ProblemDTO {
	problemDTOs := make([]ProblemDTO, 0, len(problems))
	for _, problem := range problems {
		problemDTOs = append(problemDTOs, mapProblemToDTO(&problem))
	}
	return problemDTOs
}

func mapSubmissionToDTO(submission *domain.Submission) SubmissionDTO {
	return SubmissionDTO{
		ID:     submission.ID,
		Code:   submission.Code,
		Status: string(submission.Status),
		Result: string(submission.Result),
	}
}

func mapContestToDTO(contest *domain.Contest) ContestDTO {
	return ContestDTO{
		ID:          contest.ID,
		CreatorID:   contest.Creator.ID,
		Title:       contest.Title,
		Slug:        contest.Slug,
		Description: contest.Description,
		StartsAt:    contest.StartsAt,
		EndsAt:      contest.EndsAt,
		Problems:    mapProblemsToDTOs(contest.Problems),
	}
}
