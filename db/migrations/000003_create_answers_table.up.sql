CREATE TABLE IF NOT EXISTS answers(
    id UUID NOT NULL PRIMARY KEY,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    answerer_id UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    answer TEXT NOT NULL,
    upvote INT NOT NULL DEFAULT 0,
    downvote INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
)