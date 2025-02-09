package handlers

import (
	"regexp"
)

type LogHandler interface {
	Matches(log string) (bool, error)
	StoreLog(log string) error
	GetValue() string
	GetName() string
}

type abstractLogHandler struct {
	regEx string
	name  string
}

func (lh abstractLogHandler) Matches(log string) (bool, error) {
	// Check if the string matches the pattern
	matched, err := regexp.MatchString(lh.regEx, log)
	if err != nil {
		return false, err
	}

	return matched, nil
}

func (lh abstractLogHandler) GetName() string {
	return lh.name
}

func Process(lh LogHandler, log string) error {
	// Check if the string matches the pattern
	matches, err := lh.Matches(log)
	if err != nil {
		return err
	}
	if matches {
		err = lh.StoreLog(log)
	}
	return err
}

type HandlerManager struct {
	handlers []LogHandler
}

func NewHandlerManager() HandlerManager {
	return HandlerManager{
		handlers: []LogHandler{},
	}
}

func (hm *HandlerManager) AddHandler(handler LogHandler) {
	hm.handlers = append(hm.handlers, handler)
}

func (hm *HandlerManager) ProcessLog(log string) {
	for _, lh := range hm.handlers {
		Process(lh, log)
	}
}
