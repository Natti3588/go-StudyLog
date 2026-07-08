CREATE TABLE weekly_goals (
  user_id UUID NOT NULL REFERENCES users(id),
  week_start DATE NOT NULL,
  target_min INTEGER NOT NULL CHECK (target_min > 0),
  PRIMARY KEY (user_id, week_start)
);