#!/usr/bin/env bash
# a very basic accesptance test for the webcounter app

curl -sS localhost:8000/get | grep -E '<div id="countervalue">[0-9]+</div>' || exit 1
curl -sS localhost:8000/incr | grep -E '<div id="countervalue">[0-9]+</div>' || exit 1
curl -sS localhost:8000/health | grep -E 'OK' || exit 1
curl -sS localhost:8000/reset | grep -E '<div id="countervalue">0</div>' || exit 1
curl -sS localhost:8000/metrics | grep -E '<tr><td>/metrics</td><td>1</td></tr>' || exit 1