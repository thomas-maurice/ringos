package main

import (
	"flag"
	"fmt"
	"time"

	ringo "github.com/thomas-maurice/ringos/go-client"
)

var (
	address string
	debug   bool
)

func init() {
	flag.StringVar(&address, "address", "http://10.99.4.102", "Address of the device")
	flag.BoolVar(&debug, "debug", false, "Debug mode ?")
}

func main() {
	flag.Parse()

	client, err := ringo.NewClient(&ringo.ClientConfig{
		TargetAddress: address,
	})

	if err != nil {
		panic(err)
	}

	status, err := client.Status()
	if err != nil {
		panic(err)
	}
	fmt.Println(status)

	colour, err := client.Colour()
	if err != nil {
		panic(err)
	}
	fmt.Println(colour)

	chase, err := client.Chase()
	if err != nil {
		panic(err)
	}
	fmt.Println(chase)

	breathing, err := client.Breathing()
	if err != nil {
		panic(err)
	}

	fmt.Println(breathing)

	// backs up the base state
	initialState, err := client.Colour()
	if err != nil {
		panic(err)
	}

	// makes a chase animation
	newColour := &ringo.ColourRequest{
		Mode:       "chase",
		Brightness: 200,
		Colour:     "ff0000",
	}

	_, err = client.SetColour(newColour)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	// makes a breathing animation
	newColour = &ringo.ColourRequest{
		Mode:       "breathing",
		Brightness: 200,
		Colour:     "00ff00",
	}

	_, err = client.SetColour(newColour)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	// restores the initial state
	resp, err := client.SetColour(initialState.ToColourRequest())
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
