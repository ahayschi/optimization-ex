# Optimization Examples

Simple [hill-climbing](https://en.wikipedia.org/wiki/Hill_climbing) and [simulated annealing](https://en.wikipedia.org/wiki/Simulated_annealing) techniques are implemented here to solve the example problem shown below. The objective is to minimize the tile difference between a given state and a target state. Here, the target state is the goal board.

![example img](https://github.com/ahayschi/optimization-ex/raw/master/hc-ex.png)

[Source (slide 8)](http://www.seas.upenn.edu/~cis391/Lectures/informed-search-II.pdf)

```txt
Running hill-climb...
Board Size: 3
Iterations: 5,000,000
Runtime: 7.204s
Min	Percent	Count
8	28.187	1,409,369
7	35.074	1,753,683
6	22.179	1,108,943
5	9.804	490,178
4	3.408	170,401
3	1.006	50,302
2	0.263	13,126
1	0.059	2,945
0	0.021	1,053
```

```txt
Running simulated annealing...
Board Size: 3
Iterations: 100
Runtime: 10.031s
Min	Percent	Count
8	0.000	0
7	0.000	0
6	0.000	0
5	0.000	0
4	0.000	0
3	6.000	6
2	21.000	21
1	41.000	41
0	32.000	32
```

Given enough time (and computing resources), simulated annealing should find the global minimum of 0 diff more often. This example roughly demonstrates the difference in power between the two approaches.

## Run
```sh
go run main.go
```