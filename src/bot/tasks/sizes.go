package tasks

import (
	"strings"

	"bot/utils"
)

func ParseSizes(sizes string) ([]string, error) {
	if sizes == "" || strings.ToLower(sizes) == "random" {
		return []string{"RANDOM"}, nil
	}

	rangeCheck := strings.Split(sizes, ":")

	if len(rangeCheck) > 1 {
		return utils.Frange(rangeCheck[0], rangeCheck[1])
	}

	return strings.Split(sizes, ";"), nil
}
