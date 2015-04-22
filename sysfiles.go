/*
 Playing around with Go on the Raspberry Pi, doing some trivial GPIO stuff
 using the clean but potentially slow /sys/class approach.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeOut(name string, data string) error {
	return ioutil.WriteFile(name, []byte(data), 0644)
}

func main() {
	const port = "17"
	const gpioBase = "/sys/class/gpio"
	pinfile := fmt.Sprintf(gpioBase+"/gpio%s", port)

	err := writeOut(gpioBase+"/export", port)
	// if "device or resource busy", it already exists - so just go ahead

	err = writeOut(pinfile+"/direction", "out")
	check(err)

	// led blinking
	for i := 0; i < 5; i++ {
		err = writeOut(pinfile+"/value", "0")
		check(err)
		time.Sleep(300 * time.Millisecond)

		err = writeOut(pinfile+"/value", "1")
		check(err)
		time.Sleep(300 * time.Millisecond)
	}
}
