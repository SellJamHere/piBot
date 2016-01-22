package thermo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RawTempReaderTestSuite struct {
	suite.Suite
}

func TestRawTempReaderTestSuite(t *testing.T) {
	fmt.Println("test")
	suite.Run(t, new(RawTempReaderTestSuite))
}

func givenValidRawTempLines() rawTempLines {
	return rawTempLines{
		"serial number here YES",
		"serial number here t=12345",
	}
}

func givenValidRawTempLines_noUpdate() rawTempLines {
	return rawTempLines{
		"serial number here NO",
		"serial number here t=12345",
	}
}

func givenValidRawTempLines_noT() rawTempLines {
	return rawTempLines{
		"serial number here YES",
		"serial number here 12345",
	}
}

func givenValidRawTempLines_notAnInt() rawTempLines {
	return rawTempLines{
		"serial number here NO",
		"serial number here t=abcde",
	}
}

func givenInvalidRawTempLines_empty() rawTempLines {
	return rawTempLines{}
}

func givenInvalidRawTempLines_tooMany() rawTempLines {
	return rawTempLines{"one", "two", "three"}
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_isValid_valid() {
	assert.True(suite.T(), givenValidRawTempLines().isValid())
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_isValid_invalidEmpty() {
	assert.False(suite.T(), givenInvalidRawTempLines_empty().isValid())
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_isValid_invalidTooMany() {
	assert.False(suite.T(), givenInvalidRawTempLines_tooMany().isValid())
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_hasUpdatedTemp_valid() {
	assert.True(suite.T(), givenValidRawTempLines().hasUpdatedTemp())
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_hasUpdatedTemp_invalidNoLines() {
	assert.False(suite.T(), givenInvalidRawTempLines_empty().hasUpdatedTemp())
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_hasUpdatedTemp_invalidNoYes() {
	assert.False(suite.T(), givenValidRawTempLines_noUpdate().hasUpdatedTemp())
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_temperature_valid() {
	temp, err := givenValidRawTempLines().temperature()
	assert.NoError(suite.T(), err)
	assert.InEpsilon(suite.T(), 12.345, temp.Celsius, 0.00000001)
	assert.InEpsilon(suite.T(), 54.221, temp.Fahrenheit, 0.00000001)
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_temperature_invalidNo_t() {
	temp, err := givenValidRawTempLines_noT().temperature()
	assert.Nil(suite.T(), temp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), noTError, err)
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_temperature_invalidLines() {
	temp, err := givenInvalidRawTempLines_empty().temperature()
	assert.Nil(suite.T(), temp)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), invalidLinesError, err)
}

func (suite *RawTempReaderTestSuite) TestRawTempLines_temperature_notAnInt() {
	temp, err := givenValidRawTempLines_notAnInt().temperature()
	assert.Nil(suite.T(), temp)
	assert.NotNil(suite.T(), err)
	assert.EqualError(suite.T(), err, "strconv.ParseInt: parsing \"abcde\": invalid syntax")
}
