package error

import "errors"

var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryDefault       = errors.New("category is default")
	ErrCategoryNameEmpty     = errors.New("category name is empty")
)
