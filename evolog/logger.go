package evolog

import "time"

const (
	LoggerEvolution = 0
	LoggerGeneration = 1
	LoggerEpoch = 2
	LoggerAnalysis = 3
	LoggerRScript = 4
)

type Logger struct {
	Type int `json:"message",csv:"message"`
	Message string `json:"message",csv:"message"`
	IsProgress bool `json:"isProgress"`
	Progress int `json:"progress"`
	CompleteNumber int `json:"completeNumber"`
	Timestamp time.Time `json:"time",csv:"time"`
}

func (l *Logger) NewLog(Type int, message string) *Logger {
	l.Type = Type
	l.Message = message
	l.Timestamp = time.Now()

	return l
}