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
	Count
)

var methodMap = map[string]Method{
	"Min":    Min,
	"Max":    Max,
	"Avg":    Avg,
	"Sum":    Sum,
	"Latest": Latest,
	"Count":  Count,
}

// Function to convert a string to a Method value
func MethodFromString(name string) (Method, error) {
	if method, ok := methodMap[name]; ok {
		return method, nil
	}
	return 0, fmt.Errorf("unknown method: %s", name)
}

type ThresholdMethod int

const (
	None ThresholdMethod = iota
	Eq
	Gt
	Lt
	Gte
	Lte
)

var thresholdMethodMap = map[string]ThresholdMethod{
	"None": None,
	"Eq":   Eq,
	"Gt":   Gt,
	"Lt":   Lt,
	"Gte":  Gte,
	"Lte":  Lte,
}

// Function to convert a string to a Method value
func ThresholdMethodFromString(name string) (ThresholdMethod, error) {
	if method, ok := thresholdMethodMap[name]; ok {
		return method, nil
	}
	return 0, fmt.Errorf("unknown threshold method: %s", name)
}

type NumbersHandler struct {
	abstractLogHandler
	method          Method
	numValue        float64
	count           int
	ThresholdMethod ThresholdMethod
	Threshold       float64
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
	num := 0.
	if nh.method != Count {
		// Find submatches
		matches := re.FindStringSubmatch(log)
		if len(matches) < 2 {
			return errors.New("no matching values found")
		}
		value := matches[1]

		// Convert string to float
		var err error
		num, err = strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
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
	case Count:
		nh.numValue = float64(nh.count + 1)
	}
	nh.count += 1
	return nil
}

func (nh NumbersHandler) evaluateThreshold() bool {
	switch nh.ThresholdMethod {
	case Eq:
		return nh.numValue == nh.Threshold
	case Gt:
		return nh.numValue > nh.Threshold
	case Lt:
		return nh.numValue < nh.Threshold
	case Gte:
		return nh.numValue >= nh.Threshold
	case Lte:
		return nh.numValue <= nh.Threshold
	}
	return true
}

func (nh NumbersHandler) GetValue() string {
	red := "\033[31m"
	green := "\033[32m"
	reset := "\033[0m" // Reset color to default
	prefix := ""
	suffix := ""
	if nh.ThresholdMethod != None {
		if nh.evaluateThreshold() {
			prefix = green
		} else {
			prefix = red
		}
		suffix = reset
	}
	if nh.method == Count {
		return fmt.Sprintf("%s%d%s", prefix, int(nh.numValue), suffix)
	}
	return fmt.Sprintf("%s%.2f%s", prefix, nh.numValue, suffix)
}
