#!/bin/sh

if [ -z "${AOC_SESSION}" ]; then
  echo "Must set AOC_SESSION environment variable to the session string."
  exit 1
fi

curl --cookie session=${AOC_SESSION} https://adventofcode.com/2022/leaderboard/private/view/534400.json > 534400.json
