package thermo

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	openError = "Error opening file"
)

type TempReaderTestSuite struct {
	suite.Suite
}

func TestTempReaderTestSuite(t *testing.T) {
	suite.Run(t, new(TempReaderTestSuite))
}

type fileSystemMockValid struct{}

func (fsmv *fileSystemMockValid) Open(name string) (file, error) {
	return &fileMockValid{}, nil
}

type fileMockValid struct {
	callCount int
}

func (fmv *fileMockValid) Read(p []byte) (n int, err error) {
	if fmv.callCount == 1 {
		return 0, io.EOF
	}

	writeStr := "serial number here YES\nserial number here t=12345"
	for i := 0; i < len(writeStr); i++ {
		p[i] = writeStr[i]
	}

	fmv.callCount++

	return len(writeStr), nil
}

func (fmv fileMockValid) Close() error {
	return nil
}

type fileMockInvalid struct {
	callCount int
}

func (fmv *fileMockInvalid) Read(p []byte) (n int, err error) {
	if fmv.callCount == 1 {
		return 0, io.EOF
	}

	writeStr := "serial number here NO\ninvalid"
	for i := 0; i < len(writeStr); i++ {
		p[i] = writeStr[i]
	}

	fmv.callCount++

	return len(writeStr), nil
}

func (fmv fileMockInvalid) Close() error {
	return nil
}

type fileMockInvalidSecondTime struct {
	callCount int
}

func (fmv *fileMockInvalidSecondTime) Read(p []byte) (n int, err error) {
	if fmv.callCount == 1 {
		return 0, errors.New("Error reading temp")
	}

	writeStr := "serial number here NO\ninvalid"
	for i := 0; i < len(writeStr); i++ {
		p[i] = writeStr[i]
	}

	fmv.callCount++

	return len(writeStr), nil
}

func (fmv fileMockInvalidSecondTime) Close() error {
	return nil
}

type fileSystemMockErrorOpen struct{}

func (fs *fileSystemMockErrorOpen) Open(name string) (file, error) {
	return nil, errors.New(openError)
}

type fileSystemMockErrorTemp struct{}

func (fs *fileSystemMockErrorTemp) Open(name string) (file, error) {
	return &fileMockInvalid{}, nil
}

type fileSystemMockErrorRead struct{}

func (fs *fileSystemMockErrorRead) Open(name string) (file, error) {
	return &fileMockInvalidSecondTime{}, nil
}

func givenFileSystemWithoutErrors() fileSystem {
	return new(fileSystemMockValid)
}

func givenFileSystemWithOpenError() fileSystem {
	return new(fileSystemMockErrorOpen)
}

func givenFileSystemWithReadError() fileSystem {
	return new(fileSystemMockErrorRead)
}

func givenFileSystemWithTempError() fileSystem {
	return new(fileSystemMockErrorTemp)
}

func (suite *TempReaderTestSuite) TestTempReader_readRaw_success() {
	reader := tempReader{
		serial: "serial",
		fs:     givenFileSystemWithoutErrors(),
	}

	rawTemp, err := reader.readRaw()

	assert.True(suite.T(), rawTemp.isValid())
	assert.NoError(suite.T(), err)
}

func (suite *TempReaderTestSuite) TestTempReader_readRaw_failureOpen() {
	reader := tempReader{
		serial: "serial",
		fs:     givenFileSystemWithOpenError(),
	}

	rawTemp, err := reader.readRaw()

	assert.Nil(suite.T(), rawTemp)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), openError, err.Error())
}

func (suite *TempReaderTestSuite) TestTempReader_ReadTemp_success() {
	reader := tempReader{
		serial: "serial",
		fs:     givenFileSystemWithoutErrors(),
	}

	temp, err := reader.ReadTemp()

	assert.NotNil(suite.T(), temp)
	assert.InEpsilon(suite.T(), 12.345, temp.Celsius, 0.00000001)
	assert.InEpsilon(suite.T(), 54.221, temp.Fahrenheit, 0.00000001)
	assert.NoError(suite.T(), err)
}

func (suite *TempReaderTestSuite) TestTempReader_ReadTemp_failureReadRaw() {
	reader := tempReader{
		serial: "serial",
		fs:     givenFileSystemWithOpenError(),
	}

	temp, err := reader.ReadTemp()

	assert.Nil(suite.T(), temp)
	assert.Error(suite.T(), err)
}

func (suite *TempReaderTestSuite) TestTempReader_ReadTemp_failureIsValid() {
	reader := tempReader{
		serial: "serial",
		fs:     givenFileSystemWithReadError(),
	}

	temp, err := reader.ReadTemp()

	assert.Nil(suite.T(), temp)
	assert.Error(suite.T(), err)
}

func (suite *TempReaderTestSuite) TestTempReader_ReadTemp_failureTemperature() {
	reader := tempReader{
		serial: "serial",
		fs:     givenFileSystemWithTempError(),
	}

	temp, err := reader.ReadTemp()

	assert.Nil(suite.T(), temp)
	assert.Error(suite.T(), err)
}
