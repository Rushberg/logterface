package handlers

import (
	"errors"
	"fmt"
	"logterface/utils"
	"math"
	"regexp"
	"strconv"
	"strings"
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
	diff := max - min
	finalArr := make([][]string, gh.height)
	for i := range finalArr {
		finalArr[i] = make([]string, gh.dataLength)
		for j := range finalArr[i] {
			finalArr[i][j] = " "
			if j%2 == 0 {
				finalArr[i][j] = "-"
			}
		}
	}
	prevRow := int(math.Floor(((slice[0] - min) / diff) * float64(gh.height-1)))
	// sometimes it doesn't dequeu in time hence Min
	for coulumn := 0; coulumn < utils.Min(gh.dataLength, len(slice)); coulumn++ {
		val := slice[coulumn]
		adjusted := ((val - min) / diff) * float64(gh.height-1)
		currRow := int(math.Floor(adjusted))
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

	res := ""
	for row := len(finalArr) - 1; row >= 0; row-- {
		avg_row_val := ""
		if row == 0 {
			avg_row_val = fmt.Sprintf("%.2f", min)
		} else if row == len(finalArr)-1 {
			avg_row_val = fmt.Sprintf("%.2f", max)
		} else if row%2 == 0 {
			row_range := diff / float64(gh.height)
			avg_row_val = fmt.Sprintf("%.2f", row_range*float64(row)+min+row_range/2)
		}

		res += fmt.Sprintf("%-*s │", len(fmt.Sprintf("%.2f", max)), avg_row_val)
		for column := 0; column < len(finalArr[row]); column++ {
			res += finalArr[row][column]
		}
		res += "\n"
	}

	return strings.TrimRight(res, "\n")
}
