package handlers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Method int

const (
	Min Method = iota
	Max
	Avg
	Sum
	Latest
)

var methodMap = map[string]Method{
	"Min":    Min,
	"Max":    Max,
	"Avg":    Avg,
	"Sum":    Sum,
	"Latest": Latest,
}

// Function to convert a string to a Method value
func MethodFromString(name string) (Method, error) {
	if method, ok := methodMap[name]; ok {
		return method, nil
	}
	return 0, fmt.Errorf("unknown method: %s", name)
}

type NumbersHandler struct {
	abstractLogHandler
	method   Method
	numValue float64
	count    int
}

func NewNumbersHandler(name string, regEx string, method Method) *NumbersHandler {
	nh := &NumbersHandler{
		method: method,
	}
	nh.name = name
	nh.regEx = regEx
	return nh
}

func (nh *NumbersHandler) StoreLog(log string) error {
	re := regexp.MustCompile(nh.regEx)
	// Find submatches
	matches := re.FindStringSubmatch(log)
	if len(matches) < 2 {
		return errors.New("no matching values found")
	}
	value := matches[1]

	// Convert string to float
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	if nh.count == 0 {
		nh.numValue = num
	}

	switch nh.method {
	case Min:
		if nh.numValue > num {
			nh.numValue = num
		}
	case Max:
		if nh.numValue < num {
			nh.numValue = num
		}
	case Avg:
		nh.numValue = (nh.numValue*float64(nh.count) + num) / float64(nh.count+1)
	case Sum:
		if nh.count > 0 {
			nh.numValue += num
		}
	case Latest:
		nh.numValue = num
	}
	nh.count += 1
	return nil
}

func (nh NumbersHandler) GetValue() string {
	return fmt.Sprintf("%.2f", nh.numValue)
}
