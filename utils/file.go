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
	numbers[0] = make([]float64, 256)
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
				str := currentNumber.String()
				num, err := strconv.ParseFloat(str, 64)
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				if counter == 256 {
					row++
					counter = 0
					numbers[row] = make([]float64, 256)
				}
				numbers[row][counter] = num
				// u, v := VelocityPointByFraction(float64(row), float64(counter))
				// if u != num && v != num {
				// 	fmt.Println("Read ", num, " Utrue ", u, " Vtrue ", v)
				// }
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
	}
	return nil
}
