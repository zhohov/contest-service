CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_username_length_check
        CHECK (char_length(username) BETWEEN 3 AND 32),
    CONSTRAINT users_email_length_check
        CHECK (char_length(email) <= 320),
    CONSTRAINT users_role_check
        CHECK (role IN ('user', 'admin'))
);

CREATE TABLE IF NOT EXISTS problems (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    difficulty TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT problems_title_length_check
        CHECK (char_length(title) <= 256),
    CONSTRAINT problems_slug_length_check
        CHECK (char_length(slug) <= 256),
    CONSTRAINT problems_difficulty_check
        CHECK (difficulty IN ('easy', 'medium', 'hard'))
);

CREATE TABLE IF NOT EXISTS test_cases (
    id SERIAL PRIMARY KEY,
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    input TEXT NOT NULL,
    expected TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS submissions (
    id SERIAL PRIMARY KEY,
    problem_id INTEGER NOT NULL REFERENCES problems(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    code TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    result TEXT NOT NULL DEFAULT 'unknown',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT submissions_status_check
        CHECK (status IN ('pending', 'done')),
    CONSTRAINT submissions_result_check
        CHECK (result IN ('runtime_error', 'wrong_answer', 'accepted', 'unknown'))
);

CREATE TABLE IF NOT EXISTS contests (
    id SERIAL PRIMARY KEY,
    creator_id INTEGER NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT contests_title_length_check
        CHECK (char_length(title) <= 256),
    CONSTRAINT contests_slug_length_check
        CHECK (char_length(slug) <= 256),
    CONSTRAINT contests_description_length_check
        CHECK (description IS NULL OR char_length(description) <= 1000),
    CONSTRAINT contests_time_check
        CHECK (ends_at > starts_at)
);

CREATE TABLE IF NOT EXISTS contest_problems (
    id SERIAL PRIMARY KEY,
    contest_id INTEGER NOT NULL REFERENCES contests(id) ON DELETE CASCADE,
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,

    CONSTRAINT contest_problems_unique
        UNIQUE (contest_id, problem_id)
);

CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    description TEXT,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT articles_title_length_check
        CHECK (char_length(title) <= 256),
    CONSTRAINT articles_slug_length_check
        CHECK (char_length(slug) <= 256),
    CONSTRAINT articles_description_length_check
        CHECK (description IS NULL OR char_length(description) <= 1000)
);

CREATE INDEX IF NOT EXISTS idx_test_cases_problem_id
    ON test_cases(problem_id);

CREATE INDEX IF NOT EXISTS idx_submissions_user_created_at
    ON submissions(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_submissions_problem_id
    ON submissions(problem_id);

CREATE INDEX IF NOT EXISTS idx_contest_problems_contest_id
    ON contest_problems(contest_id);

CREATE INDEX IF NOT EXISTS idx_articles_created_at
    ON articles(created_at DESC);