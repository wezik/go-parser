package parser

import "time"

type Log struct {
	Id              int
	TimestampStart  time.Time
	TimestampFinish time.Time
}

type LogTimestamp struct {
	Id        int
	State     string
	Timestamp int64 
}

const (
	StartFlag  = "STARTED"
	FinishFlag = "FINISHED"
)
