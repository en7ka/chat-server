package cmd

import (
	"time"
)

type (
	IDs []int64

	Message struct {
		ChatID  int64
		From    string
		FromUID int64
		Body    string
		Time    time.Time
	}
)
