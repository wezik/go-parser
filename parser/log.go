package parser

type Log struct {
	Id int
	TimestampStart int64
	TimestampFinish int64
}

const (
	StartFlag = "STARTED"
	FinishFlag = "FINISHED"
)
