package thermo

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Temperature struct {
	Celsius    float64
	Fahrenheit float64
}

func (t Temperature) Pretty() string {
	return fmt.Sprintf("fahrenheight: %f, celsius: %f", t.Fahrenheit, t.Celsius)
}

type TemperatureReader interface {
	Setup() error
	ReadTemp() (*Temperature, error)
}

type tempReader struct {
	serial string
}

func NewTemperatureReader(serialNumber string) TemperatureReader {
	return &tempReader{
		serial: serialNumber,
	}
}

func (t tempReader) Setup() error {
	err := exec.Command("modprobe", "w1-gpio").Run()
	if err != nil {
		return err
	}

	err = exec.Command("modprobe", "w1-therm").Run()
	if err != nil {
		return err
	}

	return nil
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

	return temp, nil
}

func (t tempReader) readRaw() (rawTempParser, error) {
	file, err := os.Open(base_dir + t.serial + file)
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
