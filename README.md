# ClockSync

This program measures a time offset from desired NTP authority and syncs a beat to it. It aligns with a daychange (still unperfect). But it gives you continuos and somewaht precise time signal to you local network via OSC, which can be of cource used many desired ways.




## Compilation

```shell
   git clone https://github.com/K0F/PerlinNet.git
   cd PerlinNet
   go mod tidy
   go build && ./PerlinNet
   ```

## Output

[![asciicast](https://asciinema.org/a/594838.svg)](https://asciinema.org/a/594838)

```
1367 0009 00013679 T 25h28m6.001s offset: -421.783µs, time: 1717630086.000524, val: 0.7686672351224568
1368 0000 00013680 T 25h28m7.001s offset: -421.783µs, time: 1717630087.001209, val: 0.6224531870009407
1368 0001 00013681 T 25h28m8.001s offset: -421.783µs, time: 1717630088.000675, val: 0.5120574887675482
1368 0002 00013682 T 25h28m9.001s offset: -421.783µs, time: 1717630089.000954, val: 0.3785339483263326
1368 0003 00013683 T 25h28m10.001s offset: -421.783µs, time: 1717630090.001417, val: 0.5004254910862553
1368 0004 00013684 T 25h28m11.001s offset: -421.783µs, time: 1717630091.001151, val: 0.7656706389200818
1368 0005 00013685 T 25h28m12.001s offset: -421.783µs, time: 1717630092.001201, val: 0.6728788622556626
```



## License

This project is licensed under the [GNU Licence](LICENSE).
