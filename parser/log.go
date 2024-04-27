package parser

import "time"

type Log struct {
	Id              int
	TimestampStart  time.Time
	TimestampFinish time.Time
}

type LogTimestamp struct {
	Id        int
	Timestamp int64 
	State     string
}

const (
	StartFlag  = "STARTED"
	FinishFlag = "FINISHED"
)
