package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func RandomString(length int, charset string) string {
	temp := make([]byte, length)

	for i := range temp {
		temp[i] = charset[rand.Intn(len(charset))]
	}

	return string(temp)
}

func Frange(r1, r2 string) ([]string, error) {

	start, err := strconv.ParseFloat(r1, 64)

	if err != nil {
		return []string{}, nil
	}

	stop, err := strconv.ParseFloat(r2, 64)

	if err != nil {
		return []string{}, nil
	}

	sizes := []string{}

	for {
		if start > stop {
			return sizes, nil
		}

		sizes = append(sizes, fmt.Sprintf("%g", start))
		start += .5
	}
}

func Extract(target, before, after string) (string, error) {
	temp := strings.Split(target, before)

	if len(temp) < 2 {
		return "", errors.New("index1 doesn't exist in target")
	}

	temp = strings.Split(temp[1], after)

	if len(temp) < 2 {
		return "", errors.New("index2 doesn't exist in target")
	}

	return temp[0], nil
}
