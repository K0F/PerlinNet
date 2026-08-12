package main

import (
	"net"
	"testing"
	"time"

	"github.com/crgimenes/go-osc"
)

func TestCalculateBeats(t *testing.T) {
	tests := []struct {
		name          string
		elapsed       time.Duration
		bpm           float64
		beatsPerBar   int
		expectedBeat  int
		expectedBar   int
		expectedTotal int
	}{
		{
			name:          "Test 1",
			elapsed:       time.Duration(0),
			bpm:           120.0,
			beatsPerBar:   4,
			expectedBeat:  0,
			expectedBar:   0,
			expectedTotal: 0,
		},
		{
			name:          "Test 2",
			elapsed:       time.Duration(time.Minute),
			bpm:           60.0,
			beatsPerBar:   4,
			expectedBeat:  0,
			expectedBar:   15,
			expectedTotal: 60,
		},
		{
			name:          "Test 3",
			elapsed:       time.Duration(time.Second * 30),
			bpm:           120.0,
			beatsPerBar:   4,
			expectedBeat:  0,
			expectedBar:   15,
			expectedTotal: 60,
		},
		{
			name:          "Test 4",
			elapsed:       time.Duration(time.Second * 15),
			bpm:           120.0,
			beatsPerBar:   4,
			expectedBeat:  2,
			expectedBar:   7,
			expectedTotal: 30,
		},
		{
			name:          "Test 5",
			elapsed:       time.Duration(time.Second * 75),
			bpm:           120.0,
			beatsPerBar:   4,
			expectedBeat:  2,
			expectedBar:   37,
			expectedTotal: 150,
		},
		{
			name:          "zero beats per bar does not panic",
			elapsed:       time.Minute,
			bpm:           60.0,
			beatsPerBar:   0,
			expectedBeat:  0,
			expectedBar:   60,
			expectedTotal: 60,
		},
		{
			name:          "negative beats per bar does not panic",
			elapsed:       time.Minute,
			bpm:           60.0,
			beatsPerBar:   -3,
			expectedBeat:  0,
			expectedBar:   60,
			expectedTotal: 60,
		},
		{
			name:          "negative bpm yields zero",
			elapsed:       time.Minute,
			bpm:           -60.0,
			beatsPerBar:   4,
			expectedBeat:  0,
			expectedBar:   0,
			expectedTotal: 0,
		},
		{
			name:          "negative elapsed yields zero",
			elapsed:       -time.Minute,
			bpm:           60.0,
			beatsPerBar:   4,
			expectedBeat:  0,
			expectedBar:   0,
			expectedTotal: 0,
		},
		{
			name:          "a whole day",
			elapsed:       24 * time.Hour,
			bpm:           120.0,
			beatsPerBar:   4,
			expectedBeat:  0,
			expectedBar:   43200,
			expectedTotal: 172800,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			beatNo, barNo, totalNo := calculateBeats(tt.elapsed, tt.bpm, tt.beatsPerBar)
			if beatNo != tt.expectedBeat || barNo != tt.expectedBar || totalNo != tt.expectedTotal {
				t.Errorf("calculateBeats(%v, %v, %v) = %v, %v, %v; want %v, %v, %v",
					tt.elapsed, tt.bpm, tt.beatsPerBar, beatNo, barNo, totalNo,
					tt.expectedBeat, tt.expectedBar, tt.expectedTotal)
			}
		})
	}
}

func TestComputeBroadcast(t *testing.T) {
	tests := []struct {
		name     string
		ip       net.IP
		mask     net.IPMask
		expected net.IP
	}{
		{
			name:     "/24",
			ip:       net.IPv4(192, 168, 0, 5),
			mask:     net.CIDRMask(24, 32),
			expected: net.IPv4(192, 168, 0, 255),
		},
		{
			name:     "/8",
			ip:       net.IPv4(10, 0, 0, 7),
			mask:     net.CIDRMask(8, 32),
			expected: net.IPv4(10, 255, 255, 255),
		},
		{
			name:     "/16",
			ip:       net.IPv4(172, 16, 5, 9),
			mask:     net.CIDRMask(16, 32),
			expected: net.IPv4(172, 16, 255, 255),
		},
		{
			name:     "/32",
			ip:       net.IPv4(127, 0, 0, 1),
			mask:     net.CIDRMask(32, 32),
			expected: net.IPv4(127, 0, 0, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeBroadcast(tt.ip.To4(), tt.mask)
			if !got.Equal(tt.expected) {
				t.Errorf("computeBroadcast(%v, %v) = %v; want %v", tt.ip, tt.mask, got, tt.expected)
			}
		})
	}
}

func TestOffsetMicros(t *testing.T) {
	tests := []struct {
		name     string
		offset   time.Duration
		expected int64
	}{
		{"negative truncates toward zero", -383364 * time.Nanosecond, -383},
		{"positive", 12500 * time.Nanosecond, 12},
		{"zero", 0, 0},
		{"exact microsecond", 1 * time.Millisecond, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := offsetMicros(tt.offset); got != tt.expected {
				t.Errorf("offsetMicros(%v) = %v; want %v", tt.offset, got, tt.expected)
			}
		})
	}
}

func TestSendToBroadcast(t *testing.T) {
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to open listener: %v", err)
	}
	defer conn.Close()

	addr := conn.LocalAddr().String()

	var client osc.Client
	msg := osc.NewMessage("/osc/timer")
	msg.Append(1717961344.0)
	msg.Append(int32(16913))
	msg.Append(int32(1))
	msg.Append(int32(140201))
	msg.Append(float32(120.0))
	msg.Append(float32(0.473502))

	expected, err := msg.MarshalBinary()
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	if err := sendToBroadcast(&client, addr, msg); err != nil {
		t.Fatalf("sendToBroadcast failed: %v", err)
	}

	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		t.Fatalf("timed out waiting for datagram: %v", err)
	}

	got := buf[:n]
	if string(got) != string(expected) {
		t.Errorf("received %q; want %q", got, expected)
	}
}

func TestClockOffsetConcurrent(t *testing.T) {
	done := make(chan struct{})
	for i := 0; i < 100; i++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for j := 0; j < 100; j++ {
				setClockOffset(time.Duration(j) * time.Microsecond)
				_ = clockOffset()
			}
		}()
	}
	for i := 0; i < 100; i++ {
		<-done
	}
}
