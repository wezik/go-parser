package parser

import "time"

type Log struct {
	Id              int
	TimestampStart  time.Time
	TimestampFinish time.Time
}

const (
	StartFlag  = "STARTED"
	FinishFlag = "FINISHED"
)
