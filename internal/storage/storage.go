package storage

import "errors"

var (
	ErrURLNotFound = errors.New("URL not foud")
	ErrURLExists   = errors.New("URL exists")
)
