CREATE TABLE user_stats (
  user_id UUID PRIMARY KEY REFERENCES users(id),
  total_min INTEGER NOT NULL DEFAULT 0,
  current_streak INTEGER NOT NULL DEFAULT 0,
  longest_streak INTEGER NOT NULL DEFAULT 0,
  last_studied_on DATE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);