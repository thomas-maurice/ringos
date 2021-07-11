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

const (
	Red   = 0
	Green = 1
	Blue  = 2
)

var (
	address         string
	debug           bool
	changeDirection bool

	// Colour random changes
	ColourChangeChance = 5
	ColourChangeMax    = 10

	// Speed random changes
	SpeedChangeChance = 5
	SpeedChangeMax    = 10

	// Direction random changes
	DirectionChangeChance = 18
	DirectionMax          = 20
)

func init() {
	host := os.Getenv("RINGOS_HOST")

	flag.StringVar(&address, "address", host, "Address of the device, or RINGOS_ADDRESS var")
	flag.BoolVar(&debug, "debug", false, "Debug mode ?")
	flag.BoolVar(&changeDirection, "chg-dir", false, "Randomly change direction ?")

	rand.Seed(time.Now().UnixNano())
}

func randomize() (string, *ringo.ColourRequest, interface{}) {
	colour := rand.Int() % 3

	var R int
	var G int
	var B int
	var chase ringo.ChaseRequest
	var clr *ringo.ColourRequest

	switch colour {
	case 0:
		R = int(rand.Int()%128) + 128
		G = rand.Int() % 255
		B = rand.Int() % 255
	case 1:
		R = rand.Int() % 255
		G = int(rand.Int()%128) + 128
		B = rand.Int() % 255
	case 2:
		R = rand.Int() % 255
		G = rand.Int() % 255
		B = int(rand.Int()%128) + 128
	}

	c := fmt.Sprintf("%02x%02x%02x", R, G, B)

	if rand.Int()%ColourChangeMax >= ColourChangeChance {
		clr = &ringo.ColourRequest{
			Colour: c,
			Mode:   ringo.MODE_CHASE,
		}
	}

	if rand.Int()%SpeedChangeMax >= SpeedChangeChance {
		chase.Speed = int64(rand.Int() % 10)
	}

	if changeDirection && rand.Int()%DirectionMax >= DirectionChangeChance {
		switch rand.Int() % 2 {
		case 1:
			chase.Direction = 1
			break
		case 0:
			chase.Direction = -1
			break
		default:
			panic("nope")
		}
	}

	return "chase", clr, &chase
}

func main() {
	flag.Parse()

	if address == "" {
		panic("ringos address cannot be empty")
	}

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
				if colour != nil {
					_, err := client.SetColour(colour)
					if err != nil {
						panic(err)
					}
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
