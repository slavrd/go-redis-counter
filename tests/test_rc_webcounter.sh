#!/usr/bin/env bash
# a very basic accesptance test for the webcounter app

curl -sS localhost:8000/get | grep -E '<div id="countervalue">[0-9]+</div>' || exit 1
curl -sS localhost:8000/incr | grep -E '<div id="countervalue">[0-9]+</div>' || exit 1