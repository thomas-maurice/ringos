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

	// Gets a client
	client, err := ringo.NewClient(&ringo.ClientConfig{
		TargetAddress: address,
	})

	if err != nil {
		panic(err)
	}

	// Lists the scanned networks, if the list is empty
	// something is likely wrong and you should retry
	scanned, err := client.ScanNetworks()
	if err != nil {
		panic(err)
	}
	fmt.Println(scanned)

	// Lists the configured networks
	known, err := client.KnownNetworks()
	if err != nil {
		panic(err)
	}
	fmt.Println(known)

	// Gets the chip's status
	status, err := client.Status()
	if err != nil {
		panic(err)
	}
	fmt.Println(status)

	// Gets the current colour config
	colour, err := client.Colour()
	if err != nil {
		panic(err)
	}
	fmt.Println(colour)

	// Gets the current chase config
	chase, err := client.Chase()
	if err != nil {
		panic(err)
	}
	fmt.Println(chase)

	// Gets the current breathing config
	breathing, err := client.Breathing()
	if err != nil {
		panic(err)
	}

	fmt.Println(breathing)

	// backs up the base colour state
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

	// changes shit a bit
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

	// makes an other breathing animation
	newBreathing := &ringo.BreathingRequest{
		BreathingResponse: ringo.BreathingResponse{
			Speed: 20,
		},
	}

	_, err = client.SetBreathing(newBreathing)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	// restores the initial state
	resp, err := client.SetColour(initialState.ToColourRequest())
	if err != nil {
		panic(err)
	}

	_, err = client.SetChase(chase.ToChaseRequest())
	if err != nil {
		panic(err)
	}

	_, err = client.SetBreathing(breathing.ToBreathingRequest())
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
