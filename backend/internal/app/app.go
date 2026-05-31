package app

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/zhohov/contest-service/backend/internal/infrastructure"
	"github.com/zhohov/contest-service/backend/internal/repository"
	"github.com/zhohov/contest-service/backend/internal/transport/event"
	httptransport "github.com/zhohov/contest-service/backend/internal/transport/http"

	"github.com/zhohov/contest-service/backend/config"
	"github.com/zhohov/contest-service/backend/internal/domain"
	repositoryInMemory "github.com/zhohov/contest-service/backend/internal/repository/in_memory"
	"github.com/zhohov/contest-service/backend/internal/usecase"
)

func NewArticleStorage() []domain.Article {
	storage := make([]domain.Article, 0, 3)

	storage = append(storage, *domain.NewArticle("Быстрая сортировка", "quick-sort", "<h2>Быстрая сортировка</h2><p>Быстрая сортировка — один из самых популярных алгоритмов сортировки, который использует стратегию «разделяй и властвуй». Алгоритм выбирает опорный элемент, после чего разбивает массив на две части: элементы меньше опорного и элементы больше опорного.</p><br><p>После разбиения каждая часть сортируется рекурсивно. В среднем quick sort работает за O(n log n), что делает его очень эффективным для больших объёмов данных. На практике этот алгоритм часто оказывается быстрее многих других сортировок благодаря хорошей локальности памяти.</p><br><p>Ключевой частью алгоритма является процедура partition. Именно она переставляет элементы массива относительно опорного значения. Существует несколько реализаций partition: Хоара, Ломуто и другие варианты.</p>"))
	storage = append(storage, *domain.NewArticle("Два указателя", "two-pointers", "<h2>Применение двух указателей</h2><p>Два указателя — это популярный алгоритмический паттерн, при котором используются два индекса для прохода по массиву или другой структуре данных. Обычно один указатель двигается с начала, а второй — с конца, либо оба движутся в одном направлении с разной скоростью.</p><br><p>Этот подход особенно полезен в задачах на отсортированные массивы. Например, с помощью двух указателей можно за O(n) находить пары чисел с заданной суммой, удалять дубликаты или искать подотрезки массива.</p><br><p>Главное преимущество метода заключается в уменьшении сложности алгоритма. Там, где наивное решение потребовало бы двойного цикла и O(n²), подход с двумя указателями позволяет получить линейное время выполнения.</p>"))
	storage = append(storage, *domain.NewArticle("Linked List", "linked-list", "<h2>Связный список</h2><p>Linked List — это структура данных, состоящая из узлов, где каждый узел хранит значение и ссылку на следующий элемент. В отличие от массива, элементы связного списка не обязаны храниться подряд в памяти.</p><br><p>Связные списки удобны в случаях, когда требуется частое добавление или удаление элементов. Например, вставка нового узла в начало списка выполняется за O(1), тогда как для массива может потребоваться перемещение элементов.</p><br><p>Существует несколько видов связанных списков: односвязные, двусвязные и циклические. Двусвязный список хранит ссылки как на следующий, так и на предыдущий элемент, что упрощает навигацию по структуре.</p>"))

	return storage
}

// func NewProblemStorage() []domain.Problem {
// 	storage := make([]domain.Problem, 0, 3)

// 	storage = append(storage, *domain.NewProblem(
// 		"Two Sum",
// 		"two-sum",
// 		"easy",
// 		"Найти два числа, сумма которых равна target",
// 		[]domain.TestCase{
// 			*domain.NewTestCase(
// 				"2 7 11 15\n9",
// 				"0 1",
// 			),

// 			*domain.NewTestCase(
// 				"3 2 4\n6",
// 				"1 2",
// 			),
// 		},
// 	))

// 	storage = append(storage, *domain.NewProblem(
// 		"Valid Parentheses",
// 		"valid-parentheses",
// 		"easy",
// 		"Проверить корректность скобочной последовательности",
// 		[]domain.TestCase{
// 			*domain.NewTestCase("()", "true"),
// 			*domain.NewTestCase("(]", "false"),
// 		},
// 	))

// 	storage = append(storage, *domain.NewProblem(
// 		"Merge Intervals",
// 		"merge-intervals",
// 		"medium",
// 		"Объединить пересекающиеся интервалы",
// 		[]domain.TestCase{
// 			*domain.NewTestCase(
// 				"1 3\n2 6\n8 10\n15 18",
// 				"[[1,6],[8,10],[15,18]]",
// 			),
// 			*domain.NewTestCase(
// 				"1 4\n4 5",
// 				"[[1,5]]",
// 			),
// 		},
// 	))

// 	return storage
// }

func NewUserStorage() []domain.User {
	storage := make([]domain.User, 0, 3)
	storage = append(storage, domain.User{
		ID:           1,
		Username:     "test",
		Email:        "test@test.com",
		Role:         domain.RoleAdmin,
		PasswordHash: "test",
	})

	storage = append(storage, domain.User{
		ID:           2,
		Username:     "test2",
		Email:        "test2@test.com",
		Role:         domain.RoleUser,
		PasswordHash: "test",
	})

	return storage
}

func Run(ctx context.Context, config *config.Config) error {
	redisClient := redis.NewClient(&redis.Options{Addr: "redis:6379", Password: "redis"})

	db, err := sql.Open("postgres", config.DB.Url)
	if err != nil {
		return fmt.Errorf("db open: %w", err)
	}
	defer db.Close()

	problemRepo := repository.NewProblemRepository(db)
	userRepo := repository.NewUserRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	contestRepo := repository.NewContestRepository(db)

	articleRepo := repositoryInMemory.NewArticleInMemoryRepo(NewArticleStorage())
	// problemRepo := repositoryInMemory.NewProblemInMemoryRepo(NewProblemStorage())
	// userRepo := repositoryInMemory.NewUserInMemoryRepo(NewUserStorage())
	jwt := infrastructure.NewJwtService([]byte("example"), time.Hour*24)
	hasher := infrastructure.NewHasher()
	broker := infrastructure.NewRedisBroker(redisClient)
	usecase := usecase.NewUseCase(articleRepo, problemRepo, userRepo, submissionRepo, contestRepo, jwt, hasher, broker)

	handler := httptransport.NewHTTPHandler(usecase, jwt)

	server := httptransport.NewServer(handler)

	httpServerErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			httpServerErr <- err
		}
	}()

	consumer := event.NewConsumer(broker, usecase)

	consumerErr := make(chan error, 1)
	go func() {
		if err := consumer.Run(ctx); err != nil {
			consumerErr <- err
		}
	}()

	select {
	case err := <-httpServerErr:
		return fmt.Errorf("http server err: %w", err)

	case err := <-consumerErr:
		return fmt.Errorf("consumer err: %w", err)

	case <-ctx.Done():
		return nil
	}
}
