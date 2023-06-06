# PerlinNet

_Perlin Noise Generator from Unix Time_

This Go program demonstrates generating Perlin noise based on Unix time. Perlin noise is a type of coherent noise that can be used for various applications, such as procedural terrain generation, procedural texture synthesis, networked live perfomances, and more!

This program uses OSC protocol to distribute value to localhost. Network distribution in next release.

You can run it on multiple places at once and as far as you have same time on machines it will produce very similar series live.


## Prerequisites

- Go programming language (version 1.2.0 or higher)

## Installation

1. Clone the repository:
   ```shell
   https://github.com/K0F/PerlinNet.git
   cd PerlinNet
   ```

## Compile

1. Run the program:
   ```shell
   go mod tidy
   go build
   ```

## Usage

1. Run the program:
   ```shell
   ./PerlinNet -p 10000 -f 60

   ```

2. The program will generate Perlin noise based on the current Unix time and display the output and send value to localhost OSC address. 
    - `-p` port on localhost to send data
    - `-f` FPS, how many messages per second there will be

## Example Output

```
time: 1685983141.799818, value: -0.291815
time: 1685983141.816417, value: -0.283964
time: 1685983141.833111, value: -0.276068
time: 1685983141.849844, value: -0.268154
time: 1685983141.866638, value: -0.260210
time: 1685983141.883528, value: -0.269776
time: 1685983141.900366, value: -0.296473
time: 1685983141.917051, value: -0.322927
time: 1685983141.933709, value: -0.349338
```

## Example SuperCollider reader

```supercollider
s.boot

p = ProxySpace.push(s)

(
    ~data.kr(1);
    ~data.mold(1);
    ~data={|x|[x].lag(1/60)};
OSCdef('/osc/perlin',{arg ... args;
	 //args.postln;
	 ~data.set(\x,args[0][2]);
},'/osc/perlin',recvPort:10000);


)

// one synth ///////////////////////////////////////

(
~one.ar(2);
~one.clock = p.clock;
~one.quant=2;
~one.fadeTime=4;
~one={
  var sig = WhiteNoise.ar(1!2);
  sig = BPF.ar(sig,~data+2*1500,0.1) * ~data;
  Splay.ar(sig,0.5,0.5);
};
~one.play;
);
~one.stop(4);
~one.clear;
```

## License

This project is licensed under the [GNU Licence](LICENSE).
