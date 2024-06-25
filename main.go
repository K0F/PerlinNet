package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"os/exec"

	//"math"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/beevik/ntp"
	"github.com/crgimenes/go-osc"

	"github.com/fatih/color"
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

	// this will not work
	client := osc.NewClient(broadcastAddr, port)
	if client == nil {
		// ... this will happen, but actually works
	}

	fmt.Printf("Starting OSC server @%v, Unix epoch: %v\n", port, time.Now().Unix())

}

func getOffset() time.Duration {

	ntpTime, err := ntp.Query("0.cz.pool.ntp.org")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("time offset from server %v\n", ntpTime.ClockOffset)
	}

	return ntpTime.ClockOffset

}

func main() {
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
	offset := time.Now().UTC().Add(ntpTime.ClockOffset) //refreshOffset(totalNo)
	beatNo, barNo, totalNo := calculateBeats(offset.Sub(midnight), *bpm, *mod)

	// Set the seed for random number generation
	//rand.New(rand.NewSource(int64(time.Now().Year())))
	//rand.Seed(int64(time.Now().Year()))

	p := perlin.NewPerlinRandSource(1.5, 2, 3, rand.NewSource(int64(time.Now().Year())))

	dur := time.Duration(60000 / *bpm) * time.Millisecond
	var drift time.Duration
	var c int = 0

	for {
		offset := time.Now().UTC().Add(ntpTime.ClockOffset) //refreshOffset(totalNo)
		t := float64(offset.UnixNano()) / 1000000000.0
		elapsed := offset.Sub(midnight)

		if elapsed > 24*time.Hour {
			start = time.Now().UTC().Add(ntpTime.ClockOffset)
			midnight = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
			offset = time.Now().Add(ntpTime.ClockOffset) //refreshOffset(totalNo)
			beatNo, barNo, totalNo = calculateBeats(offset.Sub(midnight), *bpm, *mod)
			//fmt.Printf("resetting counters %v %v %v",beatNo, barNo, totalNo)
		}

		// time.Sleep() is slightly drifting over time, correction needed here
		drift = time.Duration(elapsed.Milliseconds()%dur.Milliseconds()) * time.Millisecond

		val := p.Noise1D(t/10) + 0.5

		if beatNo == 0 {
			color.Green("%04d %04d %08d T %v offset: %v, time: %f, val: %v\n", barNo, beatNo, totalNo, elapsed.Round(time.Duration(1*time.Millisecond)), ntpTime.ClockOffset, t, val)
			if *sound {
				go runBeep("beep/sound.wav")
			}
		} else {
			fmt.Printf("%04d %04d %08d T %v offset: %v, time: %f, val: %v\n", barNo, beatNo, totalNo, elapsed.Round(time.Duration(1*time.Millisecond)), ntpTime.ClockOffset, t, val)
			//fmt.Printf("%04d %04d %08d T %v\n", barNo, beatNo, totalNo, elapsed.Round(time.Duration(1*time.Millisecond)))
		}

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
			//client.Send(msg)
			//client2.Send(msg)

		}(beatNo, totalNo, *bpm, t, val)

		totalNo = totalNo + 1
		beatNo = beatNo + 1

		if beatNo >= *mod {
			beatNo = 0
			barNo = barNo + 1
			c++
		}

		/*
			if c%10 == 0 {
				go func(_offset time.Time) {
					_offset = _offset.Add(getOffset())
					time.Sleep(1)
				}(offset)
			}
		*/

		// calculate drift correction
		ms := time.Duration(dur.Milliseconds()-drift.Milliseconds()) * time.Millisecond
		time.Sleep(ms)

		//time.Sleep(time.Duration(1000 / *fps) * time.Millisecond)
	}

}

// sendToBroadcast odesílá OSC zprávu na danou broadcast adresu
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
	return beatNo, barNo, totalBeats
}
