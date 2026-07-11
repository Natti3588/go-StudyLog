package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrCategoryNotFound   = errors.New("category not found")
	ErrStudyLogNotFound   = errors.New("study log not found")
)
