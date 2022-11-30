package errors

import "errors"

var (
	ErrOriginalURLExist = errors.New("originalURL exist")
	ErrNoSuchID         = errors.New("no such id")
)
