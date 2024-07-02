package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"os/exec"

	//	"errors"

	//"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/beevik/ntp"
	"github.com/crgimenes/go-osc"

	"github.com/fatih/color"
	//	term "github.com/nsf/termbox-go"
)

var ntpTime time.Time
var client osc.Client
var broadcastAddr string

func runBeep(arg string) {
	cmd := exec.Command("./beep/beep", arg)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error running beep: %v\n", err)
	}
}

func startServer(port int) {

	// Local broadcast adress (this will be changed if detected correctly)
	broadcastAddr = "192.168.0.255" + strconv.Itoa(port)

	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error getting network adapter:", err)
		return
	}

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			fmt.Println("Error getting interface address", iface.Name, ":", err)
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ip := ipNet.IP.To4()
				mask := ipNet.Mask
				broadcast := make(net.IP, len(ip))
				for i := range ip {
					broadcast[i] = ip[i] | ^mask[i]
				}

				fmt.Printf("Interface: %s, IP: %s, Maska: %s, Broadcast: %s\n", iface.Name, ip, mask, broadcast)
				broadcastAddr = fmt.Sprintf("%s:%s", broadcast, strconv.Itoa(port))
			}
		}
	}

	client := osc.NewClient(broadcastAddr, port)
	if client == nil {
		// ... this will happen, and/but actually works
	}

	fmt.Printf("Starting OSC server @%v, Unix epoch: %v\n", port, time.Now().Unix())

}

func getOffset() (int64, error) {

	ntpTime, err := ntp.Query("0.cz.pool.ntp.org")
	if err != nil {
		fmt.Println(err)
		return 0, err
	} else {
		color.Red("SYNC time offset from server: %v\n", ntpTime.ClockOffset)
	}

	return int64(time.Duration(ntpTime.ClockOffset * time.Microsecond)), nil

}

func main() {

	/*
		err := term.Init()
		if err != nil {
			panic(err)
		}

		defer term.Close()
	*/

	port := flag.Int("p", 10000, "Port to send OSC messages (def. 10000)")
	sound := flag.Bool("s", true, "Play a beep sound each 0nth cycle (green).")

	mod := flag.Int("m", 4, "beats per bar")
	bpm := flag.Float64("b", 120.0, "beats per minute")

	flag.Parse()

	go startServer(*port)

	ntpTime, err := ntp.Query("0.cz.pool.ntp.org")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("time offset from server %v\n", ntpTime.ClockOffset)
	}

	start := time.Now().UTC().Add(ntpTime.ClockOffset)
	midnight := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	//offset := start

	p := perlin.NewPerlinRandSource(1.5, 2, 3, rand.NewSource(int64(time.Now().Year())))

	dur := time.Duration(60000000 / *bpm) * time.Microsecond
	//var drift time.Duration
	var c int = 0

	// main loop
	for {
		offset := time.Now().UTC().Add(ntpTime.ClockOffset)
		t := float64(offset.UnixNano()) / 1000000000.0
		elapsed := offset.Sub(midnight)
		beatNo, barNo, totalNo := calculateBeats(offset.Sub(midnight), *bpm, *mod)

		if elapsed > 24*time.Hour {
			// sync to ntp server
			ntpTime, err = ntp.Query("0.cz.pool.ntp.org")
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("time offset from server %v\n", ntpTime.ClockOffset)
			}

			start = time.Now().UTC().Add(ntpTime.ClockOffset)
			midnight = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
			offset = time.Now().Add(ntpTime.ClockOffset) //refreshOffset(totalNo)
			beatNo, barNo, totalNo = calculateBeats(offset.Sub(midnight), *bpm, *mod)
			//fmt.Printf("resetting counters %v %v %v",beatNo, barNo, totalNo)
		}

		val := p.Noise1D(t/10) + 0.5

		go func(beatNo int, totalNo int, bpm float64, t float64, val float64) {
			msg := osc.NewMessage("/osc/timer")
			msg.Append(t)
			msg.Append(int32(barNo))
			msg.Append(int32(beatNo))
			msg.Append(int32(totalNo))
			msg.Append(float32(bpm))
			msg.Append(float32(val))

			// Odeslání OSC zprávy na broadcast adresu
			err = sendToBroadcast(&client, broadcastAddr, msg)
			if err != nil {
				fmt.Println("There was an error sending OSC message:", err)
			}

		}(beatNo, totalNo, *bpm, t, val)

		totalNo = totalNo + 1
		beatNo = beatNo + 1

		if beatNo >= *mod {
			beatNo = 0
			barNo = barNo + 1
			c++
		}

		if totalNo%100 == 0 {
			go func(_offset time.Time) {
				off, err := getOffset()
				if err != nil {

					_offset = time.UnixMicro(off)
					time.Sleep(1 * time.Second)
				}
			}(offset)
		}

		// time.Sleep() is slightly drifting over time, correction needed here
		drift := time.Duration(elapsed.Microseconds()%dur.Microseconds()) * time.Microsecond

		if beatNo == 0 {
			color.Green("T:%f UTC:%v OFFSET: %v VAL:%v BPM: %f BAR:%04d BEAT:%04d TOTAL:%08d\n", t, elapsed.Round(time.Duration(1*time.Millisecond)), ntpTime.ClockOffset+drift, val, *bpm, barNo, beatNo, totalNo)
			if *sound {
				go runBeep("beep/sound.wav")
			}
		} else {
			fmt.Printf("T:%f UTC:%v OFFSET: %v VAL:%v BPM: %f BAR:%04d BEAT:%04d TOTAL:%08d\n", t, elapsed.Round(time.Duration(1*time.Millisecond)), ntpTime.ClockOffset+drift, val, *bpm, barNo, beatNo, totalNo)
		}

		// calculate drift correction

		ms := time.Duration(dur.Microseconds()-drift.Microseconds()) * time.Microsecond
		time.Sleep(ms)

		//time.Sleep(time.Duration(1000 / *fps) * time.Millisecond)
	}

}

func sendToBroadcast(client *osc.Client, address string, msg *osc.Message) error {
	conn, err := net.Dial("udp", address)
	if err != nil {
		return err
	}
	defer conn.Close()

	udpConn := conn.(*net.UDPConn)
	err = udpConn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if err != nil {
		return err
	}

	data, err := msg.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = udpConn.Write(data)
	return err
}

func calculateBeats(elapsed time.Duration, bpm float64, beatsPerBar int) (int, int, int) {
	totalMinutes := elapsed.Minutes()
	totalBeats := int(totalMinutes * bpm)
	barNo := totalBeats / beatsPerBar
	beatNo := totalBeats % beatsPerBar

	// dirty way to not get -1
	if beatNo < 0 {
		beatNo = 0
	}

	if totalBeats < 0 {
		totalBeats = 0
	}

	if barNo < 0 {
		barNo = 0
	}

	return beatNo, barNo, totalBeats
}
