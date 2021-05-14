package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/filariow/garden/pkg/valve"
)

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	d := driver()
	if d == nil {
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "on":
		if err := d.SwitchOn(); err != nil {
			fmt.Println(err)
		}
	case "off":
		if err := d.SwitchOff(); err != nil {
			fmt.Println(err)
		}
	default:
		usage()
	}
}

func usage() {
	fmt.Printf(`
Set environment variables VPIN_1 and VPIN_2 to the pin number which are used to pilot the valve and run the executable

VPIN_1=7 VPIN_2=8 ./valvectl COMMAND
	- on       switch on the valve
	- off      switch off the valve
`)
}

func driver() valve.Driver {
	p1, p2 := pin(1), pin(2)
	if p1 == "" || p2 == "" || p1 == p2 {
		fmt.Printf(`Environment variables VPIN_1 ("%s") and VPIN_2 ("%s") are undefined or not valid`, p1, p2)
		usage()
		return nil
	}

	d := valve.New(p1, p2)
	return d
}

func pin(n int) string {
	e := fmt.Sprintf("VPIN_%d", n)
	p := os.Getenv(e)
	p = strings.Trim(p, " ")
	return p
}
