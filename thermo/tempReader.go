package thermo

import (
	"bufio"
	"fmt"
	"os/exec"
	"time"
)

const (
	base_dir  = "/sys/bus/w1/devices/"
	file_path = "/w1_slave"
)

type Temperature struct {
	Celsius      float64
	Fahrenheit   float64
	MeasuredTime time.Time
}

func (t Temperature) Pretty() string {
	return fmt.Sprintf("fahrenheight: %f, celsius: %f", t.Fahrenheit, t.Celsius)
}

type TemperatureReader interface {
	ReadTemp() (*Temperature, error)
}

type tempReader struct {
	serial string

	fs fileSystem
}

func NewTemperatureReader(serialNumber string) (TemperatureReader, error) {
	err := exec.Command("modprobe", "w1-gpio").Run()
	if err != nil {
		return nil, err
	}

	err = exec.Command("modprobe", "w1-therm").Run()
	if err != nil {
		return nil, err
	}

	return &tempReader{
		serial: serialNumber,
		fs:     osFS{},
	}, nil
}

func (t tempReader) ReadTemp() (*Temperature, error) {
	rawTemp, err := t.readRaw()
	if err != nil {
		return nil, err
	}

	for rawTemp.isValid() != true && rawTemp.hasUpdatedTemp() {
		time.Sleep(1 * time.Second)
		rawTemp, err = t.readRaw()
		if err != nil {
			return nil, err
		}
	}

	temp, err := rawTemp.temperature()
	if err != nil {
		return nil, err
	}

	temp.MeasuredTime = time.Now()

	return temp, nil
}

func (t tempReader) readRaw() (rawTempParser, error) {
	file, err := t.fs.Open(base_dir + t.serial + file_path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := rawTempLines{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	return lines, nil
}
