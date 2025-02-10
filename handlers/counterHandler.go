package handlers

import (
	"fmt"
)

type CounterHandler struct {
	abstractLogHandler
	count int
}

func NewCounterHandler(name string, regEx string) *CounterHandler {
	ch := &CounterHandler{}
	ch.name = name
	ch.regEx = regEx
	return ch
}

func (ch *CounterHandler) StoreLog(log string) error {
	ch.count += 1
	return nil
}

func (ch CounterHandler) GetValue() string {
	return fmt.Sprintf("%d", ch.count)
}
