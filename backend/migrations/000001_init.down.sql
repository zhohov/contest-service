DROP INDEX IF EXISTS idx_articles_created_at;
DROP INDEX IF EXISTS idx_contest_problems_contest_id;
DROP INDEX IF EXISTS idx_submissions_problem_id;
DROP INDEX IF EXISTS idx_submissions_user_created_at;
DROP INDEX IF EXISTS idx_test_cases_problem_id;

DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS contest_problems;
DROP TABLE IF EXISTS contests;
DROP TABLE IF EXISTS submissions;
DROP TABLE IF EXISTS test_cases;
DROP TABLE IF EXISTS problems;
DROP TABLE IF EXISTS users;