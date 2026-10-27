package usecase

import "github.com/zhohov/contest-service/backend/internal/domain"

func mapArticleToDTO(article *domain.Article) ArticleDTO {
	return ArticleDTO{
		ID: article.ID,
		Author: UserDTO{
			ID:       article.Author.ID,
			Username: article.Author.Username,
			Email:    article.Author.Email,
			Role:     article.Author.Role,
		},
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
		ID:        submission.ID,
		ProblemID: submission.Problem.ID,
		User: SubmissionUserDTO{
			ID:       submission.User.ID,
			Username: submission.User.Username,
		},
		Code:   submission.Code,
		Status: string(submission.Status),
		Result: string(submission.Result),
	}
}

func mapSubmissionsToDTOs(submissions []domain.Submission) []SubmissionDTO {
	submissionDTOs := make([]SubmissionDTO, 0, len(submissions))
	for _, submission := range submissions {
		submissionDTOs = append(submissionDTOs, mapSubmissionToDTO(&submission))
	}
	return submissionDTOs
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

func mapContestsToDTOs(contests []domain.Contest) []ContestDTO {
	contestDTOs := make([]ContestDTO, 0, len(contests))
	for _, contest := range contests {
		contestDTOs = append(contestDTOs, mapContestToDTO(&contest))
	}
	return contestDTOs
}
