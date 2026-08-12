# PerliNet

This program measures a time offset from desired NTP authority and syncs a beat to it.
It aligns with GMT daychange (midnight) to get the same time globally.
It gives you continuos and somewhat precise time signal to you local network via OSC,
which can be ofc used many desired ways.


![Go](https://github.com/k0f/PerlinNet/actions/workflows/go.yml/badge.svg)

## Compilation

```shell
   git clone https://github.com/K0F/PerlinNet.git
   cd PerlinNet
   make
   ./PerlinNet
   ```

`make` builds the `beep` helper (C + [miniaudio](https://miniaud.io/)) and the Go binary.

### Android (Termux)

Everything builds natively on the phone, no cross-compilation needed:

```shell
   pkg install golang clang make git
   git clone https://github.com/K0F/PerlinNet.git
   cd PerlinNet
   make termux
   ./PerlinNet
```

Notes:

- The `beep` helper is compiled for your device. If you get
  `Cannot run '.../beep': it was not built for this platform`, rebuild it with
  `make -C beep`.
- Sound requires Android 8.0+ (the miniaudio AAudio backend). If audio fails,
  you can run without sound using `-s=false`.
- NTP needs network access (UDP port 123); works on most home WiFi / mobile
  networks.

## Usage

```shell
Usage of ./PerlinNet:
  -b float
    	beats per minute (default 120)
  -m int
    	beats per bar (default 4)
  -p int
    	Port to send OSC messages (def. 10000) (default 10000)
  -s	Play a beep sound each 0nth cycle (green). (default true)
```


## OSC message (default @port 10000)

By default PerliNet detects local broadcast address and sends OSC messages to it.

It follows a pattern `/osc/timer diiiff 1717961344.000000 16913 1 140201 120.000000 0.473502`

 - `/osc/timer`, message address
 - `diiiff`, datatype pattern
 - `16913`, no. of bars
 - `1`, no. of beat (1..4) by default
 - `140201`, no. of beats (from midnight GMT), any desired BPM is recalculated
 - `120.00`, current bpm
 - `0.473502`, synchronized perlin noise value, to use for any purpose


## Fine tuning

To keep your clock in good shape consider using `chrony` or older `ntp` client. After some time of using it you will see that your system clock are nearly in perfect sync with NTP authority. 

Further precision using PTP protocol is out of scope of this software it usually requires additional hardware. Some newer network cards or recent RaspberryPIs have it on board. That are steps towards locally syncronized precision.

## Output
[![asciicast](https://asciinema.org/a/663299.svg)](https://asciinema.org/a/663299)

```
T:1717630228.976975 UTC:1h30m28.977s OFFSET: -383.364µs VAL:0.566756 BPM:120.000000 BAR:0000 BEAT:0000 TOTAL:00000000
T:1717630229.001363 UTC:1h30m29.001s OFFSET: -383.364µs VAL:0.566555 BPM:120.000000 BAR:0000 BEAT:0001 TOTAL:00000001
T:1717630230.001246 UTC:1h30m30.001s OFFSET: -383.364µs VAL:0.499974 BPM:120.000000 BAR:0000 BEAT:0002 TOTAL:00000002
T:1717630231.001023 UTC:1h30m31.001s OFFSET: -383.364µs VAL:0.558049 BPM:120.000000 BAR:0000 BEAT:0003 TOTAL:00000003
T:1717630232.000385 UTC:1h30m32.000s OFFSET: -383.364µs VAL:0.739756 BPM:120.000000 BAR:0001 BEAT:0000 TOTAL:00000004
T:1717630233.001225 UTC:1h30m33.001s OFFSET: -383.364µs VAL:0.775286 BPM:120.000000 BAR:0001 BEAT:0001 TOTAL:00000005
T:1717630234.000614 UTC:1h30m34.001s OFFSET: -383.364µs VAL:0.666320 BPM:120.000000 BAR:0001 BEAT:0002 TOTAL:00000006
T:1717630235.001155 UTC:1h30m35.001s OFFSET: -383.364µs VAL:0.575154 BPM:120.000000 BAR:0001 BEAT:0003 TOTAL:00000007
T:1717630236.000592 UTC:1h30m36.001s OFFSET: -383.364µs VAL:0.632522 BPM:120.000000 BAR:0001 BEAT:0000 TOTAL:00000008
```



## License

This project is licensed under the [GNU Licence](LICENSE).
