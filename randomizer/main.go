package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	ringo "github.com/thomas-maurice/ringos/go-client"
)

var (
	address string
	debug   bool
	Red     = 0
	Green   = 1
	Blue    = 2
)

func init() {
	flag.StringVar(&address, "address", "http://10.99.4.102", "Address of the device")
	flag.BoolVar(&debug, "debug", false, "Debug mode ?")

	rand.Seed(time.Now().UnixNano())
}

func randomize() (string, *ringo.ColourRequest, interface{}) {
	colour := rand.Int() % 3

	var R int
	var G int
	var B int

	switch colour {
	case 1:
		R = 255
		G = rand.Int() % 255
		B = rand.Int() % 255
	case 2:
		R = rand.Int() % 255
		G = 255
		B = rand.Int() % 255
	case 3:
		R = rand.Int() % 255
		G = rand.Int() % 255
		B = 255
	}

	c := fmt.Sprintf("%02x%02x%02x", R, G, B)

	fmt.Println(c)

	return "chase", &ringo.ColourRequest{
		Colour: c,
	}, nil
}

func main() {
	flag.Parse()

	// Gets a client
	client, err := ringo.NewClient(&ringo.ClientConfig{
		TargetAddress: address,
		Debug:         debug,
	})

	if err != nil {
		panic(err)
	}

	tckr := time.NewTicker(time.Duration(time.Second * 2))

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT)

	// get the base configs
	baseColour, err := client.Colour()
	if err != nil {
		panic(err)
	}
	fmt.Println(baseColour)

	// Gets the current chase config
	baseChase, err := client.Chase()
	if err != nil {
		panic(err)
	}
	fmt.Println(baseChase)

	// Gets the current breathing config
	baseBreathing, err := client.Breathing()
	if err != nil {
		panic(err)
	}
	fmt.Println(baseBreathing)

	for {
		select {
		case <-tckr.C:
			action, colour, data := randomize()
			switch action {
			case "static":
				break
			case "breathing":
				break
			case "chase":
				_, err := client.SetColour(colour)
				if err != nil {
					panic(err)
				}
				if data != nil {
					chaseData, ok := data.(*ringo.ChaseRequest)
					if !ok {
						fmt.Println("could not cast data to chase")
						return
					}
					_, err = client.SetChase(chaseData)
					if err != nil {
						panic(err)
					}
				}
			}
			break
		case <-signalChan:
			goto endOfLoop
		}
	}
endOfLoop:
	_, err = client.SetColour(baseColour.ToColourRequest())
	if err != nil {
		panic(err)
	}

	_, err = client.SetChase(baseChase.ToChaseRequest())
	if err != nil {
		panic(err)
	}

	_, err = client.SetBreathing(baseBreathing.ToBreathingRequest())
	if err != nil {
		panic(err)
	}
}
