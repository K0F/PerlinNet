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

	"github.com/fatih/color"
)

var ntpTime time.Time

func main() {
	port := flag.Int("p", 10000, "Port to send OSC messages (def. 10000)")
	//fps := flag.Int("f", 60, "Frames per second to send osc messages (def. 60)")
	//verbose := flag.Bool("v", false, "Print out values")

	mod := flag.Int("m", 4, "beats per bar")
	bpm := flag.Float64("b", 120.0, "beats per minute")

	flag.Parse()

	client := osc.NewClient("127.0.0.1", *port)

	fmt.Printf("Starting OSC server @%v, Unix epoch: %v\n", *port, time.Now().Unix())

	ntpTime, err := ntp.Query("0.cz.pool.ntp.org")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("time offset from server %v\n", ntpTime.ClockOffset)
	}
	// Set the seed for random number generation
	//rand.New(rand.NewSource(int64(time.Now().Year())))
	//rand.Seed(int64(time.Now().Year()))

	p := perlin.NewPerlinRandSource(1.5, 2, 3, rand.NewSource(int64(time.Now().Year())))

	start := time.Now().Add(ntpTime.ClockOffset)
	midnight := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	beatNo, barNo, totalNo := 0, 0, 0

	dur := time.Duration(60000 / *bpm) * time.Millisecond
	var drift time.Duration

	for {

		offset := time.Now().Add(ntpTime.ClockOffset) //refreshOffset(totalNo)
		t := float64(offset.UnixNano()) / 1000000000.0
		elapsed := offset.Sub(midnight)

		// time.Sleep() is slightly drifting over time, correction needed here
		drift = time.Duration(elapsed.Milliseconds()%dur.Milliseconds()) * time.Millisecond

		val := p.Noise1D(t/10) + 0.5

		if beatNo == 0 {
			color.Green("%04d %04d %08d T %v offset: %v, time: %f, val: %v\n", barNo, beatNo, totalNo, elapsed.Round(time.Duration(1*time.Millisecond)), ntpTime.ClockOffset, t, val)
		} else {
			fmt.Printf("%04d %04d %08d T %v offset: %v, time: %f, val: %v\n", barNo, beatNo, totalNo, elapsed.Round(time.Duration(1*time.Millisecond)), ntpTime.ClockOffset, t, val)
			//fmt.Printf("%04d %04d %08d T %v\n", barNo, beatNo, totalNo, elapsed.Round(time.Duration(1*time.Millisecond)))
		}

		go func(beatNo int, totalNo int, bpm float64, t float64, val float64) {
			msg := osc.NewMessage("/osc/timer")
			msg.Append(float64(t))
			msg.Append(int32(beatNo))
			msg.Append(int32(totalNo))
			msg.Append(float64(bpm))
			msg.Append(float64(val))
			client.Send(msg)
			//client2.Send(msg)

		}(beatNo, totalNo, *bpm, t, val)

		totalNo = totalNo + 1
		beatNo = beatNo + 1

		if beatNo >= *mod {
			beatNo = 0
			barNo = barNo + 1
		}

		// calculate drift correction
		ms := time.Duration(dur.Milliseconds()-drift.Milliseconds()) * time.Millisecond
		time.Sleep(ms)

		//time.Sleep(time.Duration(1000 / *fps) * time.Millisecond)
	}

}
