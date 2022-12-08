package main

import (
	"github.com/google/go-github/v32/github"
)

func NewTimeFormat(t github.Timestamp) string {
	return t.Format("02-Jan-2006 15:04:05")
}
