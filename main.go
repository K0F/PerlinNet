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

	p := perlin.NewPerlinRandSource(2, 2, 3, rand.NewSource(int64(time.Now().Year())))

	for {

		offset := time.Now().Add(response.ClockOffset)
		t := float64(offset.UnixNano()) / 1000000000.0
		// elapsed := t.Sub(start).Seconds()
		//val := generatePerlinNoise(int64(offset.Year()), t)
		val := p.Noise1D(t / 10)

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

/*
func generatePerlinNoise(seed int64, input float64) float64 {
	// Set the number of octaves and persistence for Perlin noise
	octaves := 4
	persistence := 0.5

	// Generate the noise value using Perlin noise algorithm
	noise := 0.0
	amplitude := 1.0
	frequency := 1.0
	for i := 0; i < octaves; i++ {
		noise += interpolateNoise(seed, input*frequency) * amplitude
		amplitude *= persistence
		frequency *= 2.0
	}

	return noise
}

func interpolateNoise(seed int64, input float64) float64 {
	// Generate random gradients for the surrounding lattice points
	grad0 := randomGradient(seed, int(math.Floor(input)))
	grad1 := randomGradient(seed, int(math.Floor(input))+1)

	// Calculate the interpolation weight
	weight := input - math.Floor(input)

	// Interpolate between the gradients
	interpolated := (1-weight)*grad0 + weight*grad1

	return interpolated
}

func randomGradient(seed int64, x int) float64 {
	// Set the seed for random number generation based on the input value
	rand.Seed(seed + int64(x))

	// Generate a random float between -1 and 1
	return -1.0 + 2.0*rand.Float64()
}
*/
