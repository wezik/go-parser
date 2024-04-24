package parser

type Log struct {
	Id int
	TimestampStart LogTimestamp
	TimestampFinish LogTimestamp
}

type LogTimestamp struct {
	Epoch int64
}

const (
	StartFlag = "STARTED"
	FinishFlag = "FINISHED"
)
