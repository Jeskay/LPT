package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func ReadDataFromFile(r io.Reader) ([][]float64, error) {
	reader := bufio.NewReader(r)
	var numbers [][]float64 = make([][]float64, 256)
	var currentNumber strings.Builder
	row := 0
	counter := 0
	for row < 256 {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if r == ' ' || r == '\r' || r == '\n' {
			if currentNumber.Len() != 0 {
				num, err := strconv.ParseFloat(currentNumber.String(), 64)
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				if counter == 256 {
					row++
					counter = 0
				}
				numbers[row] = append(numbers[row], num)
				counter++
				currentNumber.Reset()
			}
			continue
		}
		_, err = currentNumber.WriteRune(r)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return numbers, nil
}

func WriteDataToFile(data [][]float64, filename string) error {
	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	for _, l := range data {
		for _, v := range l {
			fo.WriteString(strconv.FormatFloat(v, 'e', -1, 64) + " ")
		}
		fo.WriteString("/n")
	}
	return nil
}
