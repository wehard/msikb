package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/karalabe/hid"
)

type mode int

const (
	normal  mode = 1
	gaming       = 2
	breathe      = 3
	demo         = 4
	wave         = 5
)

type region int

const (
	left       region = 1
	middle            = 2
	right             = 3
	logo              = 4
	frontLeft         = 5
	frontRight        = 6
	mouse             = 7
)

func setMode(dev *hid.Device, m mode) error {
	var buffer []byte
	buffer = make([]byte, 8)
	buffer[0] = 1
	buffer[1] = 2
	buffer[2] = 65      // Commit
	buffer[3] = byte(m) // Mode
	buffer[4] = 0
	buffer[5] = 0
	buffer[6] = 0
	buffer[7] = 236 // EOR (end of request)
	i, err := dev.SendFeatureReport(buffer)
	if err != nil {
		fmt.Println("Error when setting mode!", i)
		return err
	}
	return nil
}

func setColor(dev *hid.Device, reg region, r, b, g byte) error {
	var buffer []byte
	buffer = make([]byte, 8)
	buffer[0] = 1
	buffer[1] = 2
	buffer[3] = byte(reg) // region
	buffer[7] = 236       // end of request

	buffer[2] = 64 // rgb
	buffer[4] = r  // g
	buffer[5] = g  // r
	buffer[6] = b  // b

	i, err := dev.SendFeatureReport(buffer)
	if err != nil {
		fmt.Println("Error when setting color!", i)
		return err
	}
	return nil
}

func main() {
	var r, g, b int
	if len(os.Args) != 4 {
		fmt.Println("Usage:", os.Args[0], "r g b")
		return
	}
	r, _ = strconv.Atoi(os.Args[1])
	g, _ = strconv.Atoi(os.Args[2])
	b, _ = strconv.Atoi(os.Args[3])
	if !hid.Supported() {
		fmt.Println("error: platform not supported!")
		return
	}

	deviceInfo := hid.Enumerate(0, 0)
	if deviceInfo == nil {
		fmt.Println("Nil")
	}
	var keyboard *hid.Device
	var err error
	for _, d := range deviceInfo {
		if d.VendorID == 0x1770 && d.ProductID == 0xff00 {
			fmt.Printf("device found.\n")
			keyboard, err = d.Open()
			if err != nil {
				fmt.Printf("error: failed to open device! (vendor %#x, product %#x)\n", d.VendorID, d.ProductID)
				return
			}
			break
		}
	}
	err = setColor(keyboard, left, byte(r), byte(g), byte(b))
	err = setColor(keyboard, middle, byte(r), byte(g), byte(b))
	err = setColor(keyboard, right, byte(r), byte(g), byte(b))
	err = setMode(keyboard, normal)
	if err == nil {
		fmt.Println("colors set succesfully.")
		keyboard.Close()
		fmt.Println("closing device.")
	}
}
