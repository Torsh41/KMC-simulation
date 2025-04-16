# Simulation with Kinetic Monte Carlo method

This program is an simulation of "life" implemented with the Kinetic Monte
Carlo method with rejections
([wiki](https://en.wikipedia.org/wiki/Kinetic_Monte_Carlo)).

The simulation takes place in a bounded 2D grid, where each cell can
either be alive or dead. At the beginning, grid has a nonempty starting state,
and the system evolves throughout iterations with the following rules:

- two living cells can not be in the same location at the same time
- a living cell can move in one of 4 directions (if possible)
- an empty cell can become alive, if it has at least two neighboring living
  cells (in 4 directions only)
- a living cell that has 3 or more neighbors can die and become empty

In total, there are 9 possible events that can occur on a given cell: 4 move
events, 4 clone events and 1 death event.

## How to

Program is written and tested in `go` version `1.23.0`.

To run the program, navigate to it's directory, then go run the main.go file.

```go run main.go
```

## License

[MIT](LICENSE)

## Acknowledgments

The program uses [pixel](https://github.com/gopxl/pixel/tree/main) library for
it's graphics.
