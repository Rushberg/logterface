package handlers

import (
	"errors"
	"fmt"
	"logterface/utils"
	"math"
	"regexp"
	"strconv"
)

type GraphHandler struct {
	abstractLogHandler
	values     utils.Queue[float64]
	dataLength int
	height     int
}

func NewGraphHandler(name string, regEx string, dataLength int, height int) *GraphHandler {
	gh := &GraphHandler{}
	gh.name = name
	gh.regEx = regEx
	gh.dataLength = dataLength
	gh.height = height
	return gh
}

func (gh *GraphHandler) StoreLog(log string) error {
	re := regexp.MustCompile(gh.regEx)
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

	gh.values.Enqueue(num)
	if gh.values.Length() > gh.dataLength {
		gh.values.Dequeue()
	}
	return nil
}

func (gh *GraphHandler) GetValue() string {
	if gh.values.Length() < 2 {
		return ""
	}
	slice := gh.values.ToSlice()
	min := math.MaxFloat64
	max := -math.MaxFloat64
	for _, val := range slice {
		if val > max {
			max = val
		}
		if val < min {
			min = val
		}
	}
	height := gh.height*3 - 1
	diff := max - min
	finalArr := make([][]string, gh.height)
	for i := range finalArr {
		finalArr[i] = make([]string, gh.values.Length())
		for j := range finalArr[i] {
			finalArr[i][j] = " "
		}
	}
	prevRow := int(math.Floor(((slice[0] - min) / diff) * float64(height) / 3))
	for coulumn, val := range slice {
		adjusted := ((val - min) / diff) * float64(height)
		currRow := int(math.Floor(adjusted / 3))
		var direction int

		if currRow > prevRow {
			direction = 1
			finalArr[currRow][coulumn] = "┌"
			finalArr[prevRow][coulumn] = "┘"
		} else if currRow < prevRow {
			direction = -1
			finalArr[currRow][coulumn] = "└"
			finalArr[prevRow][coulumn] = "┐"
		} else {
			finalArr[prevRow][coulumn] = "─"
		}

		for row := prevRow + direction; math.Abs(float64(row-currRow)) > 0.1; row += direction {
			finalArr[row][coulumn] = "│"
		}

		prevRow = currRow
	}
	res := fmt.Sprintf("%.2f\n", max)
	for i := len(finalArr) - 1; i >= 0; i-- {
		for j := 0; j < len(finalArr[i]); j++ {
			res += finalArr[i][j]
		}
		res += "\n"
	}
	res += fmt.Sprintf("%.2f", min)

	return res
}
