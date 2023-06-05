# PerlinNet

_Perlin Noise Generator from Unix Time_

This Go program demonstrates generating Perlin noise based on Unix time. Perlin noise is a type of coherent noise that can be used for various applications, such as procedural terrain generation, procedural texture synthesis, networked live perfomances, and more!

This program uses OSC protocol to distribute value to localhost. Network distribution in next release.

You can run it on multiple places at one and as far as you have same time on machines it will produce completely the same series live.


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
1685968448.541762: -0.740221
```

## Contributing

Contributions are welcome! If you have any suggestions, improvements, or bug fixes, please create a pull request or open an issue.

## License

This project is licensed under the [GNU Licence](LICENSE).
