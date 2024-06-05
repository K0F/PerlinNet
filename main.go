package main

import (
	"flag"
	"fmt"
	"time"

	//"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/beevik/ntp"
	"github.com/crgimenes/go-osc"
)

func main() {
	port := flag.Int("p", 10000, "Port to send OSC messages (def. 10000)")
	fps := flag.Int("f", 60, "Frames per second to send osc messages (def. 60)")
	verbose := flag.Bool("v", false, "Print out values")

	flag.Parse()

	client := osc.NewClient("127.0.0.1", *port)

	fmt.Printf("Starting OSC server @%v, Unix epoch: %v\n", *port, time.Now().Unix())

	response, err := ntp.Query("0.cz.pool.ntp.org")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("time offset from server %v\n", response.ClockOffset)
	}
	// Set the seed for random number generation
	//rand.New(rand.NewSource(int64(time.Now().Year())))
	//rand.Seed(int64(time.Now().Year()))

	p := perlin.NewPerlinRandSource(1.5, 2, 3, rand.NewSource(int64(time.Now().Year())))

	for {

		offset := time.Now().Add(response.ClockOffset)
		t := float64(offset.UnixNano()) / 1000000000.0
		// elapsed := t.Sub(start).Seconds()
		//val := generatePerlinNoise(int64(offset.Year()), t)
		val := p.Noise1D(t/10) + 0.5

		if *verbose == true {
			fmt.Printf("offset: %v, time: %f, value: %f\n", response.ClockOffset, t, val)
		}

		go func(val float64) {
			msg := osc.NewMessage("/osc/perlin")
			msg.Append(float64(t))
			msg.Append(float64(val))
			client.Send(msg)
		}(val)

		time.Sleep(time.Duration(1000 / *fps) * time.Millisecond)
	}

}
