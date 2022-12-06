# Advent of Code 2022

[![Tests](https://github.com/devries/advent_of_code_2022/actions/workflows/main.yml/badge.svg)](https://github.com/devries/advent_of_code_2022/actions/workflows/main.yml)
[![Stars: 10](https://img.shields.io/badge/⭐_Stars-10-yellow)](https://adventofcode.com/2022)

## Plan for This Year

This year I will be using Go once again, but in an effort to do something
different I will be using [GitHub codespaces](https://docs.github.com/en/codespaces)
in order to see how well it works for someone who likes to work in the command
line. 

I'll be using the [GitHub CLI](https://cli.github.com/) in order to connect with
my codespace (using SSH) and have a [dotfiles configuration](https://github.com/devries/dotfiles)
which installs my preferred editor ([Neovim](https://neovim.io/)) and tmux so
I can work comfortably. I should be able to stick to the free tier (60 hours of
usage and 15 GB-months storage) for Advent of Code, and it gives me enough
time to see if working through codespaces is something that I enjoy. The big
down side is likely to be the delay caused by working over an SSH connection. I
am not a big fan of Neovim's remote editing feature, so I plan to run Neovim
on the remote codespace.

It seems like cloud-hosted development environments are becoming more common, so
I look forward to seeing how this experiment pans out.

## Index

- [Day 1: Calorie Counting](https://adventofcode.com/2022/day/1) - [part 1](day01_p1/main.go), [part 2](day01_p2/main.go)

  Codespaces did have a bit of a delay over ssh, though that was fairly manageable.
  When saving the file the first time it took the linter some time to get going
  and also it appears that codespaces does not have `goimports` installed by
  default, so I had to look up the standard library imports rather than rely on
  the go tooling. It's all the little things that make a development environment
  work.

- [Day 2: Rock Paper Scissors](https://adventofcode.com/2022/day/2) - [part 1](day02_p1/main.go), [part 2](day02_p2/main.go)

  I added `goimports` to my environment and things ran fairly smoothly. The delay
  is still just enough to be slightly annoying.

- [Day 3: Rucksack Reorganization](https://adventofcode.com/2022/day/3) - [part 1](day03_p1/main.go), [part 2](day03_p2/main.go)

  Bit arithmetic makes an appearance!

- [Day 4: Camp Cleanup](https://adventofcode.com/2022/day/4) - [part 1](day04_p1/main.go), [part 2](day04_p2/main.go)

- [Day 5: Supply Stacks](https://adventofcode.com/2022/day/5) - [part 1](day05_p1/main.go), [part 2](day05_p2/main.go)

  I decided to stop using codespaces. It was fine, but the small amount of delay
  over ssh was noticeable and unnecessary considering I am at my computer.

- [Day 6: Tuning Trouble](https://adventofcode.com/2022/day/6) - [part 1](day06_p1/main.go), [part 2](day06_p2/main.go)

  One thing python has which Go does not have in its standard library are tools
  to find combinations and permutations of values from a list. Luckily I wrote
  some for an earlier advent of code. I separated them into a separate library
  at [github.com/devries/combs](https://github.com/devries/combs).
