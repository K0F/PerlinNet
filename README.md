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
   go mod tidy
   go build && ./PerlinNet
   ```

## Usage

```shell
Usage of ./PerlinNet:
  -b float
    	beats per minute (default 120)
  -m int
    	beats per bar (default 4)
  -p int
    	Port to send OSC messages (def. 10000) (default 10000)
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

To keep your clock in good shape consider using `chrony` or older `ntp` client. After some time of using it you will see that your system clock are nearly in perfect sync with NTP authority. Further precision using PTP protocol usually requires additional hardware. Some newer network cards or recent RaspberryPIs have it on board. That are further steps towards near perfect precision.

## Output
[![asciicast](https://asciinema.org/a/663299.svg)](https://asciinema.org/a/663299)

```
0000 0000 00000000 T 1h30m28.977s offset: -383.364µs, time: 1717630228.976975, val: 0.5667562122176195
0000 0001 00000001 T 1h30m29.001s offset: -383.364µs, time: 1717630229.001363, val: 0.5665546624880777
0000 0002 00000002 T 1h30m30.001s offset: -383.364µs, time: 1717630230.001246, val: 0.4999737582904957
0000 0003 00000003 T 1h30m31.001s offset: -383.364µs, time: 1717630231.001023, val: 0.5580492802259492
0000 0004 00000004 T 1h30m32s offset: -383.364µs, time: 1717630232.000385, val: 0.7397558895249743
0000 0005 00000005 T 1h30m33.001s offset: -383.364µs, time: 1717630233.001225, val: 0.7752859607994234
0000 0006 00000006 T 1h30m34.001s offset: -383.364µs, time: 1717630234.000614, val: 0.6663198499634451
0000 0007 00000007 T 1h30m35.001s offset: -383.364µs, time: 1717630235.001155, val: 0.575154047879391
0000 0008 00000008 T 1h30m36.001s offset: -383.364µs, time: 1717630236.000592, val: 0.632522448253578
```



## License

This project is licensed under the [GNU Licence](LICENSE).
