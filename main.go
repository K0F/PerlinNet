package main

import (
	"flag"
	"fmt"
	"time"

	"math"
	"math/rand"

	"github.com/crgimenes/go-osc"
)

func main() {
	port := flag.Int("p", 10000, "Port to send OSC messages (def. 10000)")
	fps := flag.Int("f", 60, "Frames per second to send osc messages (def. 60)")

	flag.Parse()

	start := time.Now()

	// Set the seed for random number generation
	rand.New(rand.NewSource(int64(start.Year())))
	//rand.Seed(int64(start.Year()))

	client := osc.NewClient("127.0.0.1", *port)

	fmt.Printf("Starting OSC server @%v, Unix epoch: %v\n", *port, time.Now().Unix())

	for {
		t := (float64(time.Now().UnixNano())) / 1000000000.0
		// elapsed := t.Sub(start).Seconds()
		val := generatePerlinNoise(int64(start.Year()), t)

		fmt.Printf("time: %f, value: %f\n", t, val)

		go func(val float64) {
			msg := osc.NewMessage("/osc/perlin")
			msg.Append(float64(t))
			msg.Append(float64(val))
			client.Send(msg)
		}(val)

		time.Sleep(time.Duration(1000 / *fps) * time.Millisecond)
	}

}

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
