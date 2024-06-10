package main

import (
	"testing"
	"time"
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
