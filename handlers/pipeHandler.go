package handlers

import (
	"strings"
)

type PipeHandler struct {
	abstractLogHandler
	logs []string
}

func NewPipeHandler() *PipeHandler {
	ph := &PipeHandler{}
	ph.regEx = ".*"
	return ph
}

func (ph *PipeHandler) StoreLog(log string) error {
	ph.logs = append(ph.logs, log)
	return nil
}

func (ph *PipeHandler) GetValue() string {
	res := strings.Join(ph.logs, "\n")
	ph.logs = ph.logs[:0]
	return res
}
