# piBot
Golang sensor library for raspberry pi


## thermo
This package reads thermometer data for the [DS18B20](https://learn.adafruit.com/adafruits-raspberry-pi-lesson-11-ds18b20-temperature-sensing/overview).

```
for {
	tempReader := thermo.NewTemperatureReader("28-0000075fd190")
	temp, err := tempReader.ReadTemp()
	if err != nil {
		fmt.Println(err)
		panic("Error reading temp")
	}

	fmt.Println(temp.Pretty())
	time.Sleep(1 * time.Minute)
}
```