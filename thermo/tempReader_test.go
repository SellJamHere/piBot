package thermo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type TemperatureReaderTestSuite struct {
	suite.Suite
}

func TestTemperatureReaderTestSuite(t *testing.T) {
	fmt.Println("test")
	suite.Run(t, new(TemperatureReaderTestSuite))
}
