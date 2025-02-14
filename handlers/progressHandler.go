package handlers

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type ProgressHandler struct {
	abstractLogHandler
	numValue          float64
	DefaultTotalValue float64
	RegexTotalValue   string
	totalValue        float64
	width             int
}

func NewProgressHandler(name string, regEx string, width int) *ProgressHandler {
	nh := &ProgressHandler{
		width: width,
	}
	nh.name = name
	nh.regEx = regEx
	return nh
}

func (ph *ProgressHandler) StoreLog(log string) error {
	re := regexp.MustCompile(ph.regEx)
	ph.totalValue = ph.DefaultTotalValue
	if ph.RegexTotalValue != "" {
		re := regexp.MustCompile(ph.RegexTotalValue)
		matches := re.FindStringSubmatch(log)
		if len(matches) < 2 {
			return errors.New("no matching values found")
		}
		totalValue := matches[1]

		var err error
		ph.totalValue, err = strconv.ParseFloat(totalValue, 64)
		if err != nil {
			return err
		}
	}
	// Find progress
	matches := re.FindStringSubmatch(log)
	if len(matches) < 2 {
		return errors.New("no matching values found")
	}
	value := matches[1]

	var err error
	ph.numValue, err = strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}

	return nil
}

func (ph ProgressHandler) GetValue() string {
	prefix := "\033[34m"
	suffix := "\033[0m" // Reset color to default

	if ph.totalValue > 0 {
		bars := float64(ph.width) * (ph.numValue / ph.totalValue)
		done := int(math.Round(bars))
		remaining := ph.width - int(done)
		return fmt.Sprintf("│%s%s%s%s│", prefix, strings.Repeat("█", done), strings.Repeat(" ", remaining), suffix)
	}
	return fmt.Sprintf("│%s│", strings.Repeat(" ", ph.width))

}
