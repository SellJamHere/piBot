package main

import (
	"fmt"
	"time"

	"github.com/SellJamHere/piBot/thermo"
)

func main() {
	tempReader, err := thermo.NewTemperatureReader("28-0000075fd199")
	if err != nil {
		fmt.Println(err)
		panic("Error initializing temperature reader")
	}

	for {
		temp, err := tempReader.ReadTemp()
		if err != nil {
			fmt.Println(err)
			panic("Error reading temp")
		}

		fmt.Println(temp.Pretty())
		time.Sleep(1 * time.Minute)
	}
}
