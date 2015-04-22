/*
 Playing around with Go on the Raspberry Pi, doing some trivial GPIO stuff 
 using the unsafe/dangerous/ugly but fast mmap approach.
*/
package main

import (
	"fmt"
	"github.com/edsrzf/mmap-go"
	"os"
	"time"
	"unsafe"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	gpioBase          = 0x3F000000 + 0x200000 // values for Raspberry Pi 2 (Mod. B)
	gpioRegs          = 41                    // number of GPIO registers
	gpioRegBufferSize = gpioRegs * 4          // 32 bits hardcoded here
)

func setbit(in *uint32, bitnum uint) {
	*in |= uint32(1 << bitnum)
}

func clearbit(in *uint32, bitnum uint) {
	*in &= ^uint32(1 << bitnum)
}

func main() {

	f, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, 0644)
	check(err)
	defer f.Close()

	mapped, err := mmap.MapRegion(f, gpioRegBufferSize, mmap.RDWR, 0, gpioBase)
	check(err)

	// mapped is a byte array, let's live dangerously and "cast" this to a uint array
	regs := (*[gpioRegs]uint32)(unsafe.Pointer(&mapped[0]))

	// GPIO function select 1 is at index 1, set 3 bits for GPIO17 to output
	setbit(&regs[1], 21)
	clearbit(&regs[1], 22)
	clearbit(&regs[1], 23)

	fmt.Printf("func select register: %032b\n", regs[1])

	// now for some led blinking
	for i := 0; i < 5; i++ {
		// GPIO pin output set 0 is at index 7, output clear is at 10
		setbit(&regs[7], 17)
		time.Sleep(300 * time.Millisecond)
		setbit(&regs[10], 17)
		time.Sleep(300 * time.Millisecond)
	}
}
