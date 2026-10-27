package http

import (
	"net/http"
)

func NewServer(h *HTTPHandler) *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/articles/", h.CorsMiddleware(h.LoggingMiddleware(h.AuthMiddleware(h.CreateArticle))))
	// mux.HandleFunc("GET /api/v1/articles/", h.CorsMiddleware(h.LoggingMiddleware(h.AuthMiddleware(h.AdminOnly(h.GetArticles)))))
	mux.HandleFunc("GET /api/v1/articles/", h.CorsMiddleware(h.LoggingMiddleware((h.GetArticles))))
	// mux.HandleFunc("GET /api/v1/articles/{slug}", h.CorsMiddleware(h.LoggingMiddleware((h.GetArticleBySlug))))

	mux.HandleFunc("POST /api/v1/problems/", h.CorsMiddleware(h.LoggingMiddleware(h.CreateProblem)))
	mux.HandleFunc("GET /api/v1/problems/", h.CorsMiddleware(h.LoggingMiddleware(h.GetProblems)))
	mux.HandleFunc("GET /api/v1/problems/{slug}", h.CorsMiddleware(h.LoggingMiddleware(h.GetProblemBySlug)))

	mux.HandleFunc("POST /api/v1/submissions/", h.CorsMiddleware(h.LoggingMiddleware(h.AuthMiddleware(h.CreateSubmission))))
	mux.HandleFunc("GET /api/v1/submissions/{id}", h.CorsMiddleware(h.LoggingMiddleware(h.GetSubmissionByID)))
	mux.HandleFunc("GET /api/v1/users/me/submissions", h.CorsMiddleware(h.LoggingMiddleware(h.AuthMiddleware(h.GetSubmissionsByUserID))))

	mux.HandleFunc("POST /api/v1/auth/login/", h.CorsMiddleware(h.LoggingMiddleware(h.Login)))
	mux.HandleFunc("POST /api/v1/auth/register/", h.CorsMiddleware(h.LoggingMiddleware(h.Register)))

	mux.HandleFunc("POST /api/v1/contests", h.CorsMiddleware(h.LoggingMiddleware(h.AuthMiddleware(h.CreateContest))))
	mux.HandleFunc("POST /api/v1/contests/", h.CorsMiddleware(h.LoggingMiddleware(h.AuthMiddleware(h.CreateContest))))
	mux.HandleFunc("PUT /api/v1/contests/{id}", h.CorsMiddleware(h.LoggingMiddleware(h.AuthMiddleware(h.UpdateContest))))
	mux.HandleFunc("GET /api/v1/users/me/contests", h.CorsMiddleware(h.LoggingMiddleware(h.AuthMiddleware(h.GetContestsByUserID))))
	mux.HandleFunc("GET /api/v1/contests/{id}/submissions", h.CorsMiddleware(h.LoggingMiddleware(h.GetSubmissionsByContestID)))
	mux.HandleFunc("GET /api/v1/contests/{slug}", h.CorsMiddleware(h.LoggingMiddleware(h.GetContestBySlug)))

	handler := h.Cors(mux)

	return &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
}
