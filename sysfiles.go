/*
 Playing around with Go on the Raspberry Pi, doing some trivial GPIO stuff
 using the clean but potentially slow /sys/class approach.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"
)

const (
	gpioBase = "/sys/class/gpio"
	pinBase  = gpioBase + "/gpio"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeOut(name string, data string) error {
	return ioutil.WriteFile(name, []byte(data), 0644)
}

func pinfile(port int) string {
	return pinBase + strconv.Itoa(port)
}

func setup(port int, output bool) {
	err := writeOut(gpioBase+"/export", strconv.Itoa(port))
	// if "device or resource busy", it already exists - so just go ahead

	var value string = "in"
	if output {
		value = "out"
	}

	err = writeOut(pinfile(port)+"/direction", value)
	check(err)
}

func blink(port int) {
	setup(port, true)

	// led blinking
	for i := 0; i < 5; i++ {
		check(writeOut(pinfile(port)+"/value", "0"))
		time.Sleep(300 * time.Millisecond)

		check(writeOut(pinfile(port)+"/value", "1"))
		time.Sleep(300 * time.Millisecond)
	}
}

func getinput(port int) {
	setup(port, false)

	// check(writeOut(pinfile(port)+"/edge", "both"))

	value, err := ioutil.ReadFile(pinfile(port) + "/value")
	check(err)
	fmt.Printf("current value on port %v: %v", port, string(value))

}

func main() {
	getinput(18)
	blink(17)
}
