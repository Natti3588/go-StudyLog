package domain

import "errors"

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrStudyLogNotFound = errors.New("study log not found")
)
