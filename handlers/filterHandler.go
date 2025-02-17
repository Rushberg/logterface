package handlers

import (
	"logterface/utils"
	"regexp"
	"strings"
)

type FilterHandler struct {
	abstractLogHandler
	values     utils.Queue[string]
	dataLength int
	width      int
}

func NewFilterHandler(name string, regEx string, dataLength int, width int) *FilterHandler {
	fh := &FilterHandler{}
	fh.name = name
	fh.regEx = regEx
	fh.dataLength = dataLength
	fh.width = width
	return fh
}

func (fh *FilterHandler) StoreLog(log string) error {
	re := regexp.MustCompile(fh.regEx)
	// Find submatches
	matches := re.FindStringSubmatch(log)
	if len(matches) < 2 {
		fh.values.Enqueue(log)
	} else {
		fh.values.Enqueue(matches[1])
	}

	if fh.values.Length() > fh.dataLength {
		fh.values.Dequeue()
	}
	return nil
}

func (fh *FilterHandler) GetValue() string {
	return strings.Join(fh.values.ToSlice(), "\n")
}
