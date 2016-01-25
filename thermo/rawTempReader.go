package thermo

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	invalidLinesError = fmt.Errorf("Invalid raw temp")
	noTError          = fmt.Errorf("Unable to find t= temp")
)

type rawTempParser interface {
	isValid() bool
	hasUpdatedTemp() bool
	temperature() (*Temperature, error)
}

type rawTempLines []string

func (r rawTempLines) isValid() bool {
	return len(r) == 2
}

func (r rawTempLines) hasUpdatedTemp() bool {
	if r.isValid() != true {
		return false
	}

	firstLine := r[0]
	hasUpdatedStr := firstLine[len(firstLine)-3:]

	return hasUpdatedStr == "YES"
}

func (r rawTempLines) temperature() (*Temperature, error) {
	if r.isValid() != true {
		return nil, invalidLinesError
	}

	secondLine := r[1]
	var temp Temperature
	equalPos := strings.IndexAny(secondLine, "t=")
	if equalPos != -1 {
		tempStr := secondLine[equalPos+2:]
		tempInt, err := strconv.Atoi(tempStr)
		if err != nil {
			return nil, err
		}

		celsius := float64(tempInt) / 1000.0
		farenheit := celsius*9.0/5.0 + 32.0

		temp = Temperature{
			Celsius:    celsius,
			Fahrenheit: farenheit,
		}
	} else {
		return nil, noTError
	}

	return &temp, nil
}
